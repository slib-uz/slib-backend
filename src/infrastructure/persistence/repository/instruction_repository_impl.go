package repository

import (
	"errors"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	_errors "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type InstructionRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewInstructionRepositoryImpl(baseRepository *BaseRepository) repository.InstructionRepository {
	return &InstructionRepositoryImpl{BaseRepository: baseRepository}
}

func (this *InstructionRepositoryImpl) CreateOrUpdate(instruction *entity.InstructionEntity) error {
	model := mapper.InstructionEntityToModel(instruction)

	var existing models.InstructionModel

	err := this.db().
		Where("key = ? AND deleted_at IS NULL", model.Key).
		First(&existing).Error

	if err == nil {
		// update
		return this.db().
			Model(&existing).
			Updates(map[string]any{
				"video_link":  model.VideoLink,
				"description": model.Description,
			}).Error
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// create
		return this.db().Create(model).Error
	}

	return err
}

func (this *InstructionRepositoryImpl) Delete(id uint) error {
	return this.db().Delete(&models.InstructionModel{}, id).Error
}

func (this *InstructionRepositoryImpl) GetAll() ([]*entity.InstructionEntity, error) {
	var ms []*models.InstructionModel
	if err := this.db().Find(&ms).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	result := make([]*entity.InstructionEntity, len(ms))
	for i, m := range ms {
		result[i] = mapper.InstructionModelToEntity(m)
	}
	return result, nil
}

func (this *InstructionRepositoryImpl) GetByKey(key string) (*entity.InstructionEntity, error) {
	var model models.InstructionModel
	if err := this.db().Where("key = ?", key).First(&model).Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.InstructionModelToEntity(&model), nil
}
