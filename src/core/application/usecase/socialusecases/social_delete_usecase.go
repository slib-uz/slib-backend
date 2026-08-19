package socialusecases

import (
	"slib.uz/src/core/domain/ports/repository"
)

type SocialDestroyUseCase struct {
	repository repository.SocialRepository
}

// @inject
func NewSocialDestroyUseCase(repository repository.SocialRepository) *SocialDestroyUseCase {
	return &SocialDestroyUseCase{
		repository: repository,
	}
}

func (this *SocialDestroyUseCase) Execute(id uint) error {
	err := this.repository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
