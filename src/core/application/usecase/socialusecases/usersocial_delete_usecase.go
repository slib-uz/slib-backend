package socialusecases

import (
	"errors"

	"slib.uz/src/core/domain/ports/repository"
)

type UserSocialDestroyUseCase struct {
	repository            repository.UserSocialRepository
	userProfileRepository repository.UserProfileRepository
}

// @inject
func NewUserSocialDestroyUseCase(repository repository.UserSocialRepository, userProfileRepository repository.UserProfileRepository) *UserSocialDestroyUseCase {
	return &UserSocialDestroyUseCase{
		repository:            repository,
		userProfileRepository: userProfileRepository,
	}
}

func (this *UserSocialDestroyUseCase) Execute(id uint, userID uint) error {
	profile, err := this.userProfileRepository.GetProfileByUserID(userID)
	if err != nil {
		return err
	}

	// Проверяем, что social принадлежит пользователю
	social, err := this.repository.GetByID(id)
	if err != nil {
		return err
	}
	if social.UserProfileID != profile.ID {
		return errors.New("social does not belong to user")
	}

	err = this.repository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
