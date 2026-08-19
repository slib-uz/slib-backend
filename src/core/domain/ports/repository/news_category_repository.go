package repository

import (
	"slib.uz/src/core/domain/entity"
)

type NewsCategoryRepository interface {
	GetAll() ([]*entity.NewsCategoryEntity, error)
}
