package myapplicationusecases

import (
	"slib.uz/src/core/application/mapper"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type MyApplicationDetailUseCase struct {
	repository repository.ApplicationRepository
}

// @inject
func NewMyApplicationDetailUseCase(repository repository.ApplicationRepository) *MyApplicationDetailUseCase {
	return &MyApplicationDetailUseCase{repository: repository}
}

// Execute so'rovchining O'Z arizasini qaytaradi. Egalik GetUserAppByID ichida
// tekshiriladi (applicationId + userId birga so'raladi — begona foydalanuvchi
// boshqa ID bilan urinsa, yozuv topilmaydi). Shu sababli bu yerda
// permissionusecases.RedactApplicationContacts QO'LLANMAYDI: ariza egasi
// allaqachon "o'zi" ekani tasdiqlangan.
//
// Ma'lum va qabul qilingan chekловa (spec §3.2): javobdagi begona co-author
// (affiliation muallif) va reviewer'larning telefon/email'i redaksiyasiz
// ko'rinadi — faqat article-application (begona ariza) yo'li redaksiya
// qilindi, my-application (o'z arizasi) ataylab qamrovdan tashqarida qoldi.
func (this *MyApplicationDetailUseCase) Execute(userId, applicationId uint) (*entity2.ApplicationResponseEntity, error) {
	application, err := this.repository.GetUserAppByID(applicationId, userId)
	if err != nil {
		return nil, err
	}

	return mapper.ApplicationEntityToResponse(application), nil
}
