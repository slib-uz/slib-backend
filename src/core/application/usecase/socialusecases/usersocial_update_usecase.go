package socialusecases

import (
	"errors"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserSocialUpdateUseCase struct {
	repository            repository.UserSocialRepository
	userProfileRepository repository.UserProfileRepository
}

// @inject
func NewUserSocialUpdateUseCase(repository repository.UserSocialRepository, userProfileRepository repository.UserProfileRepository) *UserSocialUpdateUseCase {
	return &UserSocialUpdateUseCase{
		repository:            repository,
		userProfileRepository: userProfileRepository,
	}
}

func (this *UserSocialUpdateUseCase) Execute(id uint, data *entity.UserSocialInputEntity, userID uint) error {
	profile, err := this.userProfileRepository.GetProfileByUserID(userID)
	if err != nil {
		return err
	}

	social, err := this.repository.GetByID(id)
	if err != nil {
		return err
	}

	if social.UserProfileID != profile.ID {
		return errors.New("social does not belong to user")
	}

	data.UserProfileID = profile.ID
	exists, err := this.repository.Exists(profile.ID, data.SocialID)
	if err != nil {
		return err
	}
	if exists && data.SocialID != social.SocialID {
		return errors.New("social already exists")
	}
	if err := this.repository.Update(id, data); err != nil {
		return err
	}
	return nil
}
