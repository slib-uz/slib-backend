package repository

import (
	"errors"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type DraftRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewDraftRepositoryImpl(baseRepository *BaseRepository) repository.DraftRepository {
	return &DraftRepositoryImpl{BaseRepository: baseRepository}
}

func (this *DraftRepositoryImpl) Save(draft *entity.DraftEntity) error {
	var _model = mapper.DraftEntityToModel(draft)

	var oldID uint

	if err := this.db().Model(&models.DraftModel{}).Where("key = ?", draft.Key).Pluck("id", &oldID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			this.db().Create(&_model)
		}
		return err
	}
	_model.ID = oldID
	return this.db().Save(&_model).Error
}

func (this *DraftRepositoryImpl) GetByKey(key string) (*entity.DraftEntity, error) {
	var _model models.DraftModel
	if err := this.db().Where("key = ?", key).First(&_model).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.DraftModelToEntity(&_model), nil
}
