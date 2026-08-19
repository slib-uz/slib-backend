package socialusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserSocialDetailUseCase struct {
	repository repository.UserSocialRepository
}

func NewUserSocialDetailUseCase(repository repository.UserSocialRepository) *UserSocialDetailUseCase {
	return &UserSocialDetailUseCase{
		repository: repository,
	}
}

func (this *UserSocialDetailUseCase) Execute(userID uint) (*entity.UserSocialEntity, error) {
	social, err := this.repository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	return social, nil
}
