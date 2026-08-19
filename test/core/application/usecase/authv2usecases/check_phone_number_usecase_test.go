package authv2usecases_test

import (
	"errors"
	"testing"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/authv2usecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/session"
)

// fakeUserRepository — GetByPhoneNumber chaqiruvlar sonini sanaydi. Muhim
// tekshiruv shu: noto'g'ri raqamda repozitoriya umuman chaqirilmasligi
// kerak, shuning uchun faqat status emas, chaqiruvlar soni ham tekshiriladi.
// Ishlatilmagan metodlar panik qiladi.
type fakeUserRepository struct {
	getByPhoneNumberCalls int
	result                *entity.UserEntity
	err                   error
}

func (f *fakeUserRepository) GetByPhoneNumber(phoneNumber string) (*entity.UserEntity, error) {
	f.getByPhoneNumberCalls++
	return f.result, f.err
}

func (f *fakeUserRepository) GetById(id uint) (*entity.UserEntity, error)          { panic("no") }
func (f *fakeUserRepository) FindByPhone(phone string) (*entity.UserEntity, error, bool) {
	panic("no")
}
func (f *fakeUserRepository) Save(user *entity.UserEntity) (*entity.UserEntity, error) {
	panic("no")
}
func (f *fakeUserRepository) Update(id uint, user *entity.UserEntity) (*entity.UserEntity, error) {
	panic("no")
}
func (f *fakeUserRepository) GetByScienceId(scienceId string) (*entity.UserEntity, error) {
	panic("no")
}
func (f *fakeUserRepository) GetGenderStatistics() (*entity.UserGenderStatisticsEntity, error) {
	panic("no")
}
func (f *fakeUserRepository) GetAgeStatistics() (*entity.UserAgeStatisticsEntity, error) {
	panic("no")
}
func (f *fakeUserRepository) GetArticleStatisticsByYear(userID uint, year int) (*entity.UserArticleStatisticsByYearEntity, error) {
	panic("no")
}
func (f *fakeUserRepository) GetArticleStatisticsByYearRange(userID uint, fromYear, toYear int) (*entity.UserArticleStatisticsByYearRangeEntity, error) {
	panic("no")
}
func (f *fakeUserRepository) GetUserActivityStartYear(userID uint) (int, error) { panic("no") }
func (f *fakeUserRepository) GetAll(page, pageSize int, search string) (*entity.PagingEntity[entity.UserEntity], error) {
	panic("no")
}
func (f *fakeUserRepository) GetDetailByID(id uint) (*entity.UserDetailEntity, error) {
	panic("no")
}
func (f *fakeUserRepository) GetByOrcid(orcid string) (*entity.UserEntity, error) { panic("no") }
func (f *fakeUserRepository) UpdateOrcid(userID uint, orcid string) error         { panic("no") }
func (f *fakeUserRepository) CreateByPhoneNumber(tx session.Tx, user *entity.UserEntity) (uint, error) {
	panic("no")
}

// Noto'g'ri formatdagi raqamda repozitoriyaga umuman murojaat qilinmasligi
// kerak — aks holda tekshiruv DB chaqiruvidan keyin turgan bo'lardi va bu
// testsiz ko'rinmas edi.
func TestCheckPhoneNumberRejectsMalformedPhoneWithoutRepositoryCall(t *testing.T) {
	repo := &fakeUserRepository{}
	uc := authv2usecases.NewCheckPhoneNumberUseCase(repo)

	err := uc.Execute("9989012345678")

	var resp *response.Response
	if !errors.As(err, &resp) {
		t.Fatalf("*response.Response kutilgandi, %T (%v) keldi", err, err)
	}
	if resp.Status != 400 {
		t.Errorf("status 400 kutilgandi, %d keldi", resp.Status)
	}
	if repo.getByPhoneNumberCalls != 0 {
		t.Errorf("repozitoriya chaqirilmasligi kerak edi, %d marta chaqirildi", repo.getByPhoneNumberCalls)
	}
}

// To'g'ri formatdagi raqamda mavjud xulq buzilmasligi kerak: repozitoriya
// chaqiriladi.
func TestCheckPhoneNumberCallsRepositoryForWellFormedPhone(t *testing.T) {
	repo := &fakeUserRepository{err: response.NotFoundError}
	uc := authv2usecases.NewCheckPhoneNumberUseCase(repo)

	_ = uc.Execute("998901234567")

	if repo.getByPhoneNumberCalls != 1 {
		t.Errorf("repozitoriya bir marta chaqirilishi kerak edi, %d marta chaqirildi", repo.getByPhoneNumberCalls)
	}
}
