package repository

import (
	"slib.uz/src/core/domain/entity"
)

type AboutRepository interface {
	GetAll() ([]*entity.AboutEntity, error)
	GetByID(id uint) (*entity.AboutEntity, error)
}
