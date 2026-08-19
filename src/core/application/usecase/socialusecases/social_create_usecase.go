package socialusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type SocialCreateUseCase struct {
	repository repository.SocialRepository
}

// @inject
func NewSocialCreateUseCase(repository repository.SocialRepository) *SocialCreateUseCase {
	return &SocialCreateUseCase{
		repository: repository,
	}
}

func (this *SocialCreateUseCase) Execute(social *entity.SocialEntity) error {
	if err := this.repository.Create(social); err != nil {
		return err
	}
	return nil
}
