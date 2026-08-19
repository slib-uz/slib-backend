package userusecases_test

import (
	"testing"

	"slib.uz/src/core/application/usecase/userusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

// fakeUserRepo — UserRepository interfeysi embedded, faqat kerakli metod yoziladi.
type fakeUserRepo struct {
	repository.UserRepository
	calls int
}

func (f *fakeUserRepo) GetDetailByID(id uint) (*entity.UserDetailEntity, error) {
	f.calls++
	return &entity.UserDetailEntity{}, nil
}

func TestUserCanReadOwnProfile(t *testing.T) {
	repo := &fakeUserRepo{}
	uc := userusecases.NewUserDetailUseCase(repo)

	_, err := uc.Execute(&entity.UserBasicEntity{ID: 42}, 42)
	if err != nil {
		t.Fatalf("o'z profilini o'qiy olmadi: %v", err)
	}
}

func TestUserCannotReadOtherProfile(t *testing.T) {
	repo := &fakeUserRepo{}
	uc := userusecases.NewUserDetailUseCase(repo)

	_, err := uc.Execute(&entity.UserBasicEntity{ID: 42}, 99)
	if err == nil {
		t.Fatal("begona profil ochildi — IDOR")
	}
	if repo.calls != 0 {
		t.Fatalf("rad etilgan so'rov uchun DB'ga borildi: %d", repo.calls)
	}
}

func TestAdminCanReadAnyProfile(t *testing.T) {
	repo := &fakeUserRepo{}
	uc := userusecases.NewUserDetailUseCase(repo)

	admin := &entity.UserBasicEntity{
		ID:    1,
		Roles: []*entity.UserRoleEntity{{Role: enum.RoleAdmin}},
	}

	if _, err := uc.Execute(admin, 99); err != nil {
		t.Fatalf("admin begona profilni o'qiy olmadi: %v", err)
	}
}

func TestNilUserIsRejected(t *testing.T) {
	repo := &fakeUserRepo{}
	uc := userusecases.NewUserDetailUseCase(repo)

	if _, err := uc.Execute(nil, 42); err == nil {
		t.Fatal("nil foydalanuvchi qabul qilindi")
	}
}
