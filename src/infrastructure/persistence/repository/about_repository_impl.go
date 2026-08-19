package repository

import (
	"errors"
	"slib.uz/src/infrastructure/db"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type AboutRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewAboutRepository(db *db.Database) repository.AboutRepository {
	return &AboutRepositoryImpl{BaseRepository: NewBaseRepository(db)}
}

func (this *AboutRepositoryImpl) GetAll() ([]*entity.AboutEntity, error) {
	var _models []*models.AboutModel
	err := this.db().Find(&_models).Error
	if err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.AboutModelListToEntityList(_models), nil
}

func (this *AboutRepositoryImpl) GetByID(id uint) (*entity.AboutEntity, error) {
	var model models.AboutModel
	err := this.db().First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or return not found error depending on project convention, but typically nil, nil or specific error
		}
		return nil, infraError.Wrap(err)
	}
	return mapper.AboutModelToEntity(&model), nil
}
