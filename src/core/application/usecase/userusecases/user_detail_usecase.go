package userusecases

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserDetailUseCase struct {
	repository repository.UserRepository
}

// @inject
func NewUserDetailUseCase(repository repository.UserRepository) *UserDetailUseCase {
	return &UserDetailUseCase{repository: repository}
}

// Execute foydalanuvchi profilini qaytaradi.
//
// Javobda (UserDetailEntity) na PINFL, na pasport, na tug'ilgan sana, na email,
// na telefon bor — u {id, F.I.SH., science ID, rasm, roles, maqolalar soni} va
// ish joyi, ilmiy ko'rsatkichlar, ilmiy daraja/unvon, ORCID dan iborat.
// Cheklovning sababi roles massivi: u foydalanuvchi qaysi jurnal va
// nashriyotlarda qanday lavozimda ekanini ochadi, ya'ni tashkiliy tuzilmani
// istalgan hisob ID bo'yicha yig'ib olardi. Shuning uchun profil faqat
// egasining o'ziga va adminga ochiq. Ommaviy muallif ma'lumoti uchun
// /api/author/... va /api/users/find endpointlari mavjud.
func (this *UserDetailUseCase) Execute(user *entity2.UserBasicEntity, id uint) (*entity2.UserDetailEntity, error) {
	if user == nil {
		return nil, response.UnauthorizedError
	}

	if user.ID != id && !permissionusecases.IsAdmin(user) {
		return nil, response.PermissionDeniedError
	}

	return this.repository.GetDetailByID(id)
}
