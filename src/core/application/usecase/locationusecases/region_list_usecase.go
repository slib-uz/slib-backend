package locationusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type RegionListUseCase struct {
	repository repository.RegionRepository
}

// @inject
func NewRegionListUseCase(repository repository.RegionRepository) *RegionListUseCase {
	return &RegionListUseCase{repository: repository}
}

func (this *RegionListUseCase) Execute() ([]*entity.RegionEntity, error) {
	return this.repository.List()
}
