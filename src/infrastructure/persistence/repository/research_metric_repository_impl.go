package repository

import (
	"errors"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type ResearchMetricRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewResearchMetricRepository(baseRepository *BaseRepository) repository.ResearchMetricRepository {
	return &ResearchMetricRepositoryImpl{BaseRepository: baseRepository}
}

func (this *ResearchMetricRepositoryImpl) UpdateOrCreate(userID uint, source enum.ResearchMetricEnum, entity *entity.ResearchMetricEntity) error {
	var existing *models.ResearchMetricModel
	err := this.db().Unscoped().Where("user_id = ? AND source = ?", userID, source).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_model := mapper.ResearchMetricEntityToModel(entity)
			_model.UserID = userID
			_model.Source = source
			if createErr := this.db().Create(&_model).Error; createErr != nil {
				return createErr
			}
		} else {
			return err
		}
	} else {
		existing.HIndex = entity.HIndex
		existing.ProfileUrl = entity.ProfileUrl
		existing.Source = source

		if existing.DeletedAt.Valid {
			existing.DeletedAt = gorm.DeletedAt{}
			if updateErr := this.db().Unscoped().Save(&existing).Error; updateErr != nil {
				return updateErr
			}
		} else {
			if updateErr := this.db().Save(&existing).Error; updateErr != nil {
				return updateErr
			}
		}
	}
	return nil
}

func (this *ResearchMetricRepositoryImpl) GetByUserID(userID uint) ([]*entity.ResearchMetricEntity, error) {
	var researchMetrics []*models.ResearchMetricModel
	if err := this.db().Where("user_id = ?", userID).Find(&researchMetrics).Error; err != nil {
		return nil, err
	}

	return mapper.ResearchMetricListModelToEntity(researchMetrics), nil
}

func (this *ResearchMetricRepositoryImpl) GetByUserIDs(userIDs []uint) (map[uint][]*entity.ResearchMetricEntity, error) {
	if len(userIDs) == 0 {
		return make(map[uint][]*entity.ResearchMetricEntity), nil
	}

	var researchMetrics []*models.ResearchMetricModel
	if err := this.db().Where("user_id IN ?", userIDs).Find(&researchMetrics).Error; err != nil {
		return nil, err
	}

	result := make(map[uint][]*entity.ResearchMetricEntity)
	for _, metric := range researchMetrics {
		if result[metric.UserID] == nil {
			result[metric.UserID] = make([]*entity.ResearchMetricEntity, 0)
		}
		result[metric.UserID] = append(result[metric.UserID], mapper.ResearchMetricModelToEntity(metric))
	}

	return result, nil
}

func (this *ResearchMetricRepositoryImpl) DeleteByIDAndUserID(id, userID uint) error {
	result := this.db().Where("id = ? AND user_id = ?", id, userID).Delete(&models.ResearchMetricModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
