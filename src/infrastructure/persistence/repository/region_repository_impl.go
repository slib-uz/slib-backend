package repository

import (
	"slib.uz/src/core/domain/entity"
	domainrepo "slib.uz/src/core/domain/ports/repository"
	_errors "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type RegionRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewRegionRepositoryImpl(base *BaseRepository) domainrepo.RegionRepository {
	return &RegionRepositoryImpl{BaseRepository: base}
}

func (r *RegionRepositoryImpl) List() ([]*entity.RegionEntity, error) {
	var ms []models.RegionModel
	if err := r.db().Order("soato asc").Find(&ms).Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.RegionModelsToEntities(ms), nil
}
