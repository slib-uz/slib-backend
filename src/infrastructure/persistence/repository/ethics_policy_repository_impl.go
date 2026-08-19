package repository

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type EthicsPolicyRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewEthicsPolicyRepository(baseRepository *BaseRepository) repository.EthicsPolicyRepository {
	return &EthicsPolicyRepositoryImpl{BaseRepository: baseRepository}
}

func (this *EthicsPolicyRepositoryImpl) GetByPaging(page, pageSize int) (*entity.PagingEntity[entity.EthicsPolicyEntity], error) {
	var total int64
	var ethicsPolicyModels []*models.EthicsPolicyModel

	countQuery := this.db().Model(&models.EthicsPolicyModel{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	query := this.db().Model(&models.EthicsPolicyModel{})
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&ethicsPolicyModels).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	ethicsPolicies := mapper.EthicsPolicyModelListToEntityList(ethicsPolicyModels)

	return entity.NewPagingEntity(page, pageSize, total, ethicsPolicies), nil
}

func (this *EthicsPolicyRepositoryImpl) GetByID(id uint) (*entity.EthicsPolicyEntity, error) {
	var ethicsPolicyModel models.EthicsPolicyModel

	query := this.db().Model(&models.EthicsPolicyModel{})
	if err := query.Where("id = ?", id).Find(&ethicsPolicyModel).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	ethicsPolicy := mapper.EthicsPolicyModelToEntity(&ethicsPolicyModel)

	return ethicsPolicy, nil
}
