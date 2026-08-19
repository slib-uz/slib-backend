package repository

import (
	"fmt"

	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type PartnerRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewPartnerRepository(repository *BaseRepository) repository.PartnerRepository {
	return &PartnerRepositoryImpl{BaseRepository: repository}
}

func (this PartnerRepositoryImpl) GetAll(page, pageSize int) (*entity2.PagingEntity[entity2.PartnerEntity], error) {
	var total int64
	var _model []*models.PartnerModel

	countQuery := this.db().Model(&models.PartnerModel{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("error counting partners: %w", err)
	}

	query := this.db().Model(&models.PartnerModel{})
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&_model).Error; err != nil {
		return nil, fmt.Errorf("error fetching partners: %w", err)
	}

	partners := mapper.PartnerModelListToEntityList(_model)

	return entity2.NewPagingEntity(page, pageSize, total, partners), nil
}
