package repository

import (
	"slib.uz/src/core/domain/entity"
)

type ClientRepository interface {
	GetByClientID(clientID string) (*entity.ClientEntity, error)
	GetById(id uint) (*entity.ClientEntity, error)
}
