package repository

import (
	"slib.uz/src/core/domain/entity"
)

type DoctorateRepository interface {
	GetByUserID(userID uint) (*entity.DoctorateEntity, error)
	BulkCreate(userID uint, doctorates []*entity.DoctorateEntity) error
}
