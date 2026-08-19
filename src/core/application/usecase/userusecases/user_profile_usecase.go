package userusecases

import (
	"errors"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserProfileUseCase struct {
	repository               repository.UserProfileRepository
	jobRepository            repository.JobRepository
	researchMetricRepository repository.ResearchMetricRepository
	articleRepository        repository.ArticleRepository
	applicationRepository    repository.ApplicationRepository
}

// @inject
func NewUserProfileUseCase(repository repository.UserProfileRepository, jobRepository repository.JobRepository, researchMetricRepository repository.ResearchMetricRepository, articleRepository repository.ArticleRepository, applicationRepository repository.ApplicationRepository) *UserProfileUseCase {
	return &UserProfileUseCase{
		repository:               repository,
		jobRepository:            jobRepository,
		researchMetricRepository: researchMetricRepository,
		articleRepository:        articleRepository,
		applicationRepository:    applicationRepository,
	}
}

// Execute foydalanuvchi profilini qaytaradi.
//
// Bitta usecase ikkita yo'nalishga xizmat qiladi: `/user/profile` (o'ziniki) va
// `/user/profile/{id}` (boshqasiniki). Ikkinchisi ommaviy muallif kartochkasini
// to'ldiradi, shuning uchun u yopilmaydi — lekin bog'lanish ma'lumotlari
// (email, telefon) begonaga ko'rinmasligi kerak. Shu sababli so'rovchi
// foydalanuvchi ham qabul qilinadi: o'zi va admin to'liq profilni oladi,
// qolganlar uchun email va telefon bo'shatiladi. Javobning boshqa qismi
// (F.I.SH., ORCID, science ID, ish joyi, ilmiy ko'rsatkichlar, bio, shahar)
// har ikkala holatda ham o'zgarmaydi.
func (this *UserProfileUseCase) Execute(
	user *entity2.UserBasicEntity, userID uint,
) (*entity2.UserProfileEntity, error) {
	profile, err := this.repository.GetUserProfileByUserID(userID)
	if err != nil {
		if errors.Is(err, response.NotFoundError) {
			return entity2.NewUserProfileEntity(0, 0, "", "", "", nil, nil, "", nil, nil, nil, nil, nil, nil, nil, 0, 0), nil
		}
		return nil, err
	}

	jobs, err := this.jobRepository.GetByUserID(userID)
	if err != nil || len(jobs) == 0 {
		profile.Job = nil
	} else {
		profile.Job = jobs[0]
	}

	researchMetrics, err := this.researchMetricRepository.GetByUserID(userID)
	if err != nil {
		profile.ResearchMetrics = []*entity2.ResearchMetricEntity{}
	} else {
		profile.ResearchMetrics = researchMetrics
	}

	publishedArticlesCount, err := this.articleRepository.GetArticleCountByUserID(userID, true)
	if err != nil {
		publishedArticlesCount = 0
	}

	processingArticlesCount, err := this.applicationRepository.GetApplicationCountByUserID(userID, false)
	if err != nil {
		processingArticlesCount = 0
	}

	profile.PublishedArticlesCount = publishedArticlesCount
	profile.ProcessingArticlesCount = processingArticlesCount

	// Egasi va admin to'liq profilni ko'radi; boshqa hamma uchun bog'lanish
	// ma'lumotlari olib tashlanadi. Foydalanuvchi noma'lum bo'lsa ham
	// olib tashlanadi — shubha bo'lsa yopiq tomonga.
	if user == nil || (user.ID != userID && !permissionusecases.IsAdmin(user)) {
		profile.Email = ""
		profile.Phone = ""
	}

	return profile, nil
}
