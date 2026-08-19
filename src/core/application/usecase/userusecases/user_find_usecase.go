package userusecases

import (
	"errors"

	"slib.uz/src/core/application/mapper"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
)

type UserFindUseCase struct {
	repository repository.UserRepository
	gateway    gateway.ScienceIDGateway
}

// @inject
func NewUserFindUseCase(repository repository.UserRepository, gateway gateway.ScienceIDGateway) *UserFindUseCase {
	return &UserFindUseCase{repository: repository, gateway: gateway}
}

// Execute berilgan scienceId bo'yicha foydalanuvchini qaytaradi. So'rovchi
// (requester) o'zi bo'lmasa yoki admin bo'lmasa, natijadagi telefon raqami
// bo'shatiladi (CWE-200: istalgan login qilgan foydalanuvchi scienceId'ni
// aylantirib boshqalarning telefonini yig'ib olmasligi uchun).
func (this *UserFindUseCase) Execute(requester *entity.UserBasicEntity, scienceId string) (*entity.UserBasicEntity, error) {
	user, err := this.repository.GetByScienceId(scienceId)

	if err != nil && !errors.Is(err, response.NotFoundError) {
		return nil, err
	}

	if user == nil {
		user, err = this.gateway.GetUserByScienceID(scienceId)
		if err != nil {
			return nil, err
		}
		user, err = this.repository.Save(user)
		if err != nil {
			return nil, err
		}
	}

	result := mapper.UserEntityToBasic(user)

	var requesterID uint
	if requester != nil {
		requesterID = requester.ID
	}
	permissionusecases.RedactBasicContact(result, requesterID, permissionusecases.IsAdmin(requester))

	return result, nil
}
