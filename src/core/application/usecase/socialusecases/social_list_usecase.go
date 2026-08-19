package socialusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type SocialAllUseCase struct {
	repository repository.SocialRepository
}

// @inject
func NewSocialAllUseCase(repository repository.SocialRepository) *SocialAllUseCase {
	return &SocialAllUseCase{
		repository: repository,
	}
}

func (this *SocialAllUseCase) Execute() ([]*entity.SocialEntity, error) {
	socials, err := this.repository.GetAll()
	if err != nil {
		return nil, err
	}
	return socials, nil
}
