package socialusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type SocialUpdateUseCase struct {
	repository repository.SocialRepository
}

// @inject
func NewSocialUpdateUseCase(repository repository.SocialRepository) *SocialUpdateUseCase {
	return &SocialUpdateUseCase{
		repository: repository,
	}
}

func (this *SocialUpdateUseCase) Execute(id uint, data *entity.SocialEntity) error {
	if err := this.repository.Update(id, data); err != nil {
		return err
	}
	return nil
}
