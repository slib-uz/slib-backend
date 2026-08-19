package socialusecases

import (
	"errors"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserSocialCreateUseCase struct {
	repository            repository.UserSocialRepository
	userProfileRepository repository.UserProfileRepository
}

// @inject
func NewUserSocialCreateUseCase(repository repository.UserSocialRepository, userProfileRepository repository.UserProfileRepository) *UserSocialCreateUseCase {
	return &UserSocialCreateUseCase{
		repository:            repository,
		userProfileRepository: userProfileRepository,
	}
}

func (this *UserSocialCreateUseCase) Execute(social *entity.UserSocialInputEntity, userID uint) error {
	profile, err := this.userProfileRepository.GetProfileByUserID(userID)
	if err != nil {
		return err
	}
	exists, err := this.repository.Exists(profile.ID, social.SocialID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("social already exists")
	}
	social.UserProfileID = profile.ID
	if err := this.repository.Create(social); err != nil {
		return err
	}
	return nil
}
