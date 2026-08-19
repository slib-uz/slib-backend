package guideusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type GuideListUseCase struct {
	repository repository.GuideRepository
}

// @inject
func NewGuideListUseCase(repository repository.GuideRepository) *GuideListUseCase {
	return &GuideListUseCase{repository: repository}
}

func (this *GuideListUseCase) Execute(page, pageSize int) (*entity2.PagingEntity[entity2.GuideListEntity], error) {
	paging, err := this.repository.GetAll(page, pageSize)

	if err != nil {
		return nil, err
	}

	return paging, nil
}
