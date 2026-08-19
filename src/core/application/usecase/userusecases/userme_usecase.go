package userusecases

import (
	"errors"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserMeUseCase struct {
	repository repository.UserProfileRepository
}

// @inject
func NewUserMeUseCase(repository repository.UserProfileRepository) *UserMeUseCase {
	return &UserMeUseCase{repository: repository}
}

func (this *UserMeUseCase) Execute(userID uint) (*entity.UserMeEntity, error) {

	user, err := this.repository.GetByUserID(userID)
	if err != nil {

		if errors.Is(err, response.NotFoundError) {
			return entity.NewUserMeEntity(0, "", "", nil, nil, 0), err
		}

		return nil, err
	}

	return user, nil
}
