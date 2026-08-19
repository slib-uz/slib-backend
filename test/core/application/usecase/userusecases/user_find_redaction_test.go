package userusecases_test

import (
	"testing"

	"slib.uz/src/core/application/usecase/userusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

// Topilgan foydalanuvchining ID'si — testlarda "o'zi" va "begona"ni ajratish uchun.
const foundUserID uint = 7

type fakeUserFindRepo struct {
	repository.UserRepository
	user *entity.UserEntity
}

func (f *fakeUserFindRepo) GetByScienceId(scienceId string) (*entity.UserEntity, error) {
	return f.user, nil
}

// fakeUserFindGateway - ScienceIDGateway metodlaridan hech biri chaqirilmasligi
// kerak, chunki repository allaqachon foydalanuvchini topadi. Chaqirilsa, panika
// bilan testni yiqitadi.
type fakeUserFindGateway struct{}

func (f *fakeUserFindGateway) GetAuthorByScienceID(scienceID string) (*entity.AuthorEntity, error) {
	panic("GetAuthorByScienceID chaqirilmasligi kerak edi")
}

func (f *fakeUserFindGateway) GetUserScientificDataByScienceID(scienceID string) (*entity.AcademicDegreeEntity, *entity.AcademicTitleEntity, error) {
	panic("GetUserScientificDataByScienceID chaqirilmasligi kerak edi")
}

func (f *fakeUserFindGateway) GetUserByScienceID(scienceID string) (*entity.UserEntity, error) {
	panic("GetUserByScienceID chaqirilmasligi kerak edi")
}

func newUserFindUseCase(user *entity.UserEntity) *userusecases.UserFindUseCase {
	return userusecases.NewUserFindUseCase(
		&fakeUserFindRepo{user: user},
		&fakeUserFindGateway{},
	)
}

func buildFoundUser() *entity.UserEntity {
	return &entity.UserEntity{
		ID:          foundUserID,
		ScienceID:   "SCI-7",
		FullName:    "Topilgan Foydalanuvchi",
		PhoneNumber: "+998901112233",
	}
}

// Asosiy talab (CWE-200): /users/find?scienceId=X begona (login qilgan, lekin
// egasi ham admin ham bo'lmagan) so'rovchiga topilgan foydalanuvchining
// telefon raqamini ko'rsatmasligi kerak.
func TestUserFindRedactsPhoneForStranger(t *testing.T) {
	uc := newUserFindUseCase(buildFoundUser())

	requester := &entity.UserBasicEntity{ID: 999}
	result, err := uc.Execute(requester, "SCI-7")
	if err != nil {
		t.Fatalf("kutilmagan xato: %v", err)
	}
	if result.PhoneNumber != "" {
		t.Errorf("begonaga telefon ochildi: %q", result.PhoneNumber)
	}
}

func TestUserFindKeepsPhoneForOwner(t *testing.T) {
	uc := newUserFindUseCase(buildFoundUser())

	requester := &entity.UserBasicEntity{ID: foundUserID}
	result, err := uc.Execute(requester, "SCI-7")
	if err != nil {
		t.Fatalf("kutilmagan xato: %v", err)
	}
	if result.PhoneNumber == "" {
		t.Error("o'ziga o'z telefoni ko'rinishi kerak edi")
	}
}

func TestUserFindKeepsPhoneForAdmin(t *testing.T) {
	uc := newUserFindUseCase(buildFoundUser())

	requester := &entity.UserBasicEntity{ID: 999, IsAdmin: true}
	result, err := uc.Execute(requester, "SCI-7")
	if err != nil {
		t.Fatalf("kutilmagan xato: %v", err)
	}
	if result.PhoneNumber == "" {
		t.Error("admin uchun telefon ko'rinishi kerak edi")
	}
}

// nil so'rovchi (masalan noto'g'ri kontekst) xavfsiz tomonga: hech kim ega
// bo'lmaydi, telefon bo'shatiladi, panika bo'lmaydi.
func TestUserFindNilRequesterRedactsPhone(t *testing.T) {
	uc := newUserFindUseCase(buildFoundUser())

	result, err := uc.Execute(nil, "SCI-7")
	if err != nil {
		t.Fatalf("kutilmagan xato: %v", err)
	}
	if result.PhoneNumber != "" {
		t.Errorf("nil so'rovchiga telefon ochildi: %q", result.PhoneNumber)
	}
}
