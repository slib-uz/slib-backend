package repository

import (
	"slib.uz/src/core/domain/entity"
	domainrepo "slib.uz/src/core/domain/ports/repository"
	_errors "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type DistrictRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewDistrictRepositoryImpl(base *BaseRepository) domainrepo.DistrictRepository {
	return &DistrictRepositoryImpl{BaseRepository: base}
}

func (r *DistrictRepositoryImpl) ListByRegionID(regionID uint) ([]*entity.DistrictEntity, error) {
	var ms []models.DistrictModel
	if err := r.db().Where("region_id = ?", regionID).Order("soato asc").Find(&ms).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	return mapper.DistrictModelsToEntities(ms), nil
}
