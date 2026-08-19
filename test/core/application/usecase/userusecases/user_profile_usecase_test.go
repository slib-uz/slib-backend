package userusecases_test

import (
	"testing"

	"slib.uz/src/core/application/usecase/userusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

// Profil egasining ID'si — testlarda "o'zi" va "begona" ni ajratish uchun.
const profileOwnerID uint = 42

type fakeUserProfileRepo struct {
	repository.UserProfileRepository
}

func (f *fakeUserProfileRepo) GetUserProfileByUserID(userID uint) (*entity.UserProfileEntity, error) {
	bio := "biografiya"
	city := enum.CityCode("TSH")
	orcid := "0000-0002-1825-0097"
	return &entity.UserProfileEntity{
		ID:        1,
		UserID:    userID,
		FullName:  "Norqulov Muhriddin",
		ScienceID: "SCI-1",
		Email:     "user@example.com",
		Phone:     "+998901234567",
		Bio:       &bio,
		City:      &city,
		ORCIDID:   &orcid,
	}, nil
}

type fakeJobRepo struct {
	repository.JobRepository
}

func (f *fakeJobRepo) GetByUserID(userID uint) ([]*entity.JobEntity, error) {
	return []*entity.JobEntity{{ID: 5}}, nil
}

type fakeResearchMetricRepo struct {
	repository.ResearchMetricRepository
}

func (f *fakeResearchMetricRepo) GetByUserID(userID uint) ([]*entity.ResearchMetricEntity, error) {
	return []*entity.ResearchMetricEntity{{ID: 9}}, nil
}

type fakeArticleRepo struct {
	repository.ArticleRepository
}

func (f *fakeArticleRepo) GetArticleCountByUserID(userID uint, isPublished bool) (int64, error) {
	return 3, nil
}

type fakeApplicationRepo struct {
	repository.ApplicationRepository
}

func (f *fakeApplicationRepo) GetApplicationCountByUserID(userID uint, isPublished bool) (int64, error) {
	return 2, nil
}

func newProfileUseCase() *userusecases.UserProfileUseCase {
	return userusecases.NewUserProfileUseCase(
		&fakeUserProfileRepo{},
		&fakeJobRepo{},
		&fakeResearchMetricRepo{},
		&fakeArticleRepo{},
		&fakeApplicationRepo{},
	)
}

// assertPayloadIntact — bog'lanish maydonlaridan tashqari hamma narsa
// har ikkala ko'rinishda ham o'zgarishsiz qolishi kerak.
func assertPayloadIntact(t *testing.T, profile *entity.UserProfileEntity) {
	t.Helper()

	if profile.FullName != "Norqulov Muhriddin" {
		t.Fatalf("F.I.SH. o'zgarib ketdi: %q", profile.FullName)
	}
	if profile.ScienceID != "SCI-1" {
		t.Fatalf("science ID o'zgarib ketdi: %q", profile.ScienceID)
	}
	if profile.Bio == nil || *profile.Bio != "biografiya" {
		t.Fatal("bio yo'qoldi")
	}
	if profile.City == nil {
		t.Fatal("shahar yo'qoldi")
	}
	if profile.ORCIDID == nil {
		t.Fatal("ORCID yo'qoldi")
	}
	if profile.Job == nil {
		t.Fatal("ish joyi yo'qoldi")
	}
	if len(profile.ResearchMetrics) != 1 {
		t.Fatalf("ilmiy ko'rsatkichlar yo'qoldi: %d", len(profile.ResearchMetrics))
	}
	if profile.PublishedArticlesCount != 3 || profile.ProcessingArticlesCount != 2 {
		t.Fatalf(
			"maqola hisoblagichlari o'zgarib ketdi: %d / %d",
			profile.PublishedArticlesCount, profile.ProcessingArticlesCount,
		)
	}
}

// Asosiy talab: /user/profile/{id} ommaviy muallif kartochkasi uchun ochiq
// qoladi, lekin begona foydalanuvchining emaili va telefoni ko'rinmaydi.
func TestStrangerDoesNotSeeContactFields(t *testing.T) {
	uc := newProfileUseCase()

	profile, err := uc.Execute(&entity.UserBasicEntity{ID: 99}, profileOwnerID)
	if err != nil {
		t.Fatalf("begona profil o'qilmadi (endpoint yopilib qolgan?): %v", err)
	}
	if profile.Email != "" {
		t.Fatalf("begonaga email ochildi: %q", profile.Email)
	}
	if profile.Phone != "" {
		t.Fatalf("begonaga telefon ochildi: %q", profile.Phone)
	}
	assertPayloadIntact(t, profile)
}

func TestOwnerSeesOwnContactFields(t *testing.T) {
	uc := newProfileUseCase()

	profile, err := uc.Execute(&entity.UserBasicEntity{ID: profileOwnerID}, profileOwnerID)
	if err != nil {
		t.Fatalf("o'z profilini o'qiy olmadi: %v", err)
	}
	if profile.Email != "user@example.com" {
		t.Fatalf("foydalanuvchi o'z emailini ko'ra olmadi: %q", profile.Email)
	}
	if profile.Phone != "+998901234567" {
		t.Fatalf("foydalanuvchi o'z telefonini ko'ra olmadi: %q", profile.Phone)
	}
	assertPayloadIntact(t, profile)
}

func TestAdminSeesContactFields(t *testing.T) {
	uc := newProfileUseCase()

	admin := &entity.UserBasicEntity{ID: 1, IsAdmin: true}

	profile, err := uc.Execute(admin, profileOwnerID)
	if err != nil {
		t.Fatalf("admin profilni o'qiy olmadi: %v", err)
	}
	if profile.Email == "" || profile.Phone == "" {
		t.Fatalf("admin uchun bog'lanish maydonlari o'chirildi: %q / %q", profile.Email, profile.Phone)
	}
}

// RoleAdmin tarmog'i ham bayroq bilan bir xil natija berishi kerak.
func TestRoleAdminSeesContactFields(t *testing.T) {
	uc := newProfileUseCase()

	admin := &entity.UserBasicEntity{
		ID:    1,
		Roles: []*entity.UserRoleEntity{{Role: enum.RoleAdmin}},
	}

	profile, err := uc.Execute(admin, profileOwnerID)
	if err != nil {
		t.Fatalf("RoleAdmin profilni o'qiy olmadi: %v", err)
	}
	if profile.Email == "" || profile.Phone == "" {
		t.Fatalf("RoleAdmin uchun bog'lanish maydonlari o'chirildi: %q / %q", profile.Email, profile.Phone)
	}
}

// Foydalanuvchi noma'lum bo'lsa yopiq tomonga: panika ham, sizib chiqish ham bo'lmasin.
func TestNilUserDoesNotSeeContactFields(t *testing.T) {
	uc := newProfileUseCase()

	profile, err := uc.Execute(nil, profileOwnerID)
	if err != nil {
		t.Fatalf("nil foydalanuvchi uchun kutilmagan xato: %v", err)
	}
	if profile.Email != "" || profile.Phone != "" {
		t.Fatalf("nil foydalanuvchiga bog'lanish maydonlari ochildi: %q / %q", profile.Email, profile.Phone)
	}
}
