package socialusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserSocialUseCase struct {
	repository repository.UserSocialRepository
}

func NewUserSocialUseCase(repository repository.UserSocialRepository) *UserSocialUseCase {
	return &UserSocialUseCase{
		repository: repository,
	}
}

func (this *UserSocialUseCase) Execute() ([]*entity.UserSocialEntity, error) {
	socials, err := this.repository.GetAll()
	if err != nil {
		return nil, err
	}
	return socials, nil
}
