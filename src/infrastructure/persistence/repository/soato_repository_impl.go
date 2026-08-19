package repository

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/db"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type SoatoRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewSoatoRepository(db *db.Database) repository.SoatoRepository {
	return &SoatoRepositoryImpl{BaseRepository: NewBaseRepository(db)}
}

func (this *SoatoRepositoryImpl) GetAll(page, size int, isChildren bool) (*entity.PagingEntity[entity.SoatoClassifierEntity], error) {
	var _models []*models.SoatoClassifierModel
	query := this.db().Model(&_models).Order("created_at desc")

	if isChildren {
		query = query.Where("parent_id IS NOT NULL").Preload("Parent")
	} else {
		query = query.Where("parent_id IS NULL").Preload("Children")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	if err := query.Limit(size).Offset((page - 1) * size).Find(&_models).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	entities := mapper.SoatoModelListToEntityList(_models)
	return entity.NewPagingEntity(page, size, total, entities), nil
}
