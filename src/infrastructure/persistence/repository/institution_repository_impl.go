package repository

import (
	entity2 "slib.uz/src/core/domain/entity"

	"gorm.io/gorm"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type InstitutionRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewInstitutionRepository(baseRepository *BaseRepository) repository.InstitutionRepository {
	return &InstitutionRepositoryImpl{BaseRepository: baseRepository}
}

func (this *InstitutionRepositoryImpl) GetByID(id uint) (*entity2.InstitutionEntity, error) {
	var institution models.InstitutionModel
	if err := this.db().Preload("Publishers").First(&institution, id).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.InstitutionModelToDetailEntity(&institution), nil
}

func (this *InstitutionRepositoryImpl) GetList(page, size int, tin, name *string) (*entity2.PagingEntity[entity2.InstitutionEntity], error) {
	var institutions []*models.InstitutionModel

	query := this.db().Model(&models.InstitutionModel{})

	if tin != nil && *tin != "" {
		query = query.Where("tin ILIKE ?", "%"+*tin+"%")
	}
	if name != nil && *name != "" {
		query = query.Where("name ILIKE ?", "%"+*name+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	if err := query.Limit(size).Offset((page - 1) * size).Order("created_at DESC").Find(&institutions).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	entities := make([]*entity2.InstitutionEntity, len(institutions))
	for i, institution := range institutions {
		entities[i] = mapper.InstitutionModelToEntity(institution)
	}

	return entity2.NewPagingEntity(page, size, total, entities), nil
}

func (this *InstitutionRepositoryImpl) Create(e *entity2.InstitutionEntity) error {
	model := mapper.InstitutionEntityToModel(e)
	if err := this.db().Create(model).Error; err != nil {
		return infraError.Wrap(err)
	}
	*e = *mapper.InstitutionModelToEntity(model)
	return nil
}

func (this *InstitutionRepositoryImpl) Update(id uint, e *entity2.InstitutionEntity) error {
	var model models.InstitutionModel
	if err := this.db().First(&model, id).Error; err != nil {
		return infraError.Wrap(err)
	}

	newModel := mapper.InstitutionEntityToModel(e)
	model.Name = newModel.Name
	model.Tin = newModel.Tin
	model.Logo = newModel.Logo

	if err := this.db().Save(&model).Error; err != nil {
		return infraError.Wrap(err)
	}
	*e = *mapper.InstitutionModelToEntity(&model)
	return nil
}

func (this *InstitutionRepositoryImpl) Delete(id uint) error {
	if err := this.db().Delete(&models.InstitutionModel{}, id).Error; err != nil {
		return infraError.Wrap(err)
	}
	return nil
}

func (this *InstitutionRepositoryImpl) SetPublishers(institutionID uint, publisherIDs []uint) error {
	return this.db().Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&models.InstitutionModel{}, institutionID).Error; err != nil {
			return infraError.Wrap(err)
		}

		if len(publisherIDs) == 0 {
			return nil
		}

		var count int64
		if err := tx.Model(&models.PublisherModel{}).
			Where("id IN ?", publisherIDs).
			Count(&count).Error; err != nil {
			return infraError.Wrap(err)
		}
		if count != int64(len(publisherIDs)) {
			return response.InvalidArgument
		}

		if err := tx.Model(&models.PublisherModel{}).
			Where("id IN ?", publisherIDs).
			Update("institution_id", institutionID).Error; err != nil {
			return infraError.Wrap(err)
		}

		return nil
	})
}

func (this *InstitutionRepositoryImpl) DetachPublishers(institutionID uint, publisherIDs []uint) (int64, error) {
	if err := this.db().First(&models.InstitutionModel{}, institutionID).Error; err != nil {
		return 0, infraError.Wrap(err)
	}

	result := this.db().Model(&models.PublisherModel{}).
		Where("id IN ? AND institution_id = ?", publisherIDs, institutionID).
		Update("institution_id", nil)

	if result.Error != nil {
		return 0, infraError.Wrap(result.Error)
	}

	return result.RowsAffected, nil
}
