package repository

import (
	"slib.uz/src/core/domain/entity"
)

type DegreeRepository interface {
	GetByUserID(userID uint) (*entity.DegreeEntity, error)
	UpdateOrCreate(userID uint, degree *entity.DegreeEntity) error
}
