package locationusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type DistrictListUseCase struct {
	repository repository.DistrictRepository
}

// @inject
func NewDistrictListUseCase(repository repository.DistrictRepository) *DistrictListUseCase {
	return &DistrictListUseCase{repository: repository}
}

func (this *DistrictListUseCase) Execute(regionID uint) ([]*entity.DistrictEntity, error) {
	return this.repository.ListByRegionID(regionID)
}
