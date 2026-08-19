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

type ClientRepositoryImpl struct {
	db *gorm.DB
}

// @inject
func NewClientRepository(db *gorm.DB) repository.ClientRepository {
	return &ClientRepositoryImpl{db: db}
}

func (r *ClientRepositoryImpl) GetByClientID(clientID string) (*entity.ClientEntity, error) {
	var model models.ClientModel
	err := r.db.Where("client_id = ? AND is_active = ?", clientID, true).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, _errors.Wrap(err)
	}

	return mapper.ClientModelToEntity(&model), nil
}

func (r *ClientRepositoryImpl) GetById(id uint) (*entity.ClientEntity, error) {
	var model models.ClientModel
	err := r.db.Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, _errors.Wrap(err)
	}

	return mapper.ClientModelToEntity(&model), nil
}
