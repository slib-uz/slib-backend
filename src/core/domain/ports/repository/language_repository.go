package repository

import (
	"slib.uz/src/core/domain/entity"
)

type LanguageRepository interface {
	GetAll() ([]*entity.LanguageEntity, error)
	IsExist(id uint) (bool, error)
	ExistingIds(ids []uint) ([]uint, error)
}
