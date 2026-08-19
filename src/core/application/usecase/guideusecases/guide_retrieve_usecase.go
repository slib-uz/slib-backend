package guideusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type GuideRetrieveUseCase struct {
	repository repository.GuideRepository
}

// @inject
func NewGuideRetrieveUseCase(repository repository.GuideRepository) *GuideRetrieveUseCase {
	return &GuideRetrieveUseCase{repository: repository}
}

func (this *GuideRetrieveUseCase) Execute(guideId uint) (*entity.GuideRetrieveEntity, error) {
	guide, err := this.repository.GetByID(guideId)

	if err != nil {
		return nil, err
	}

	return guide, nil
}
