package repository

import (
	"slib.uz/src/core/domain/entity"
)

type AcademicTitleRepository interface {
	UpdateOrCreate(academicTitle *entity.AcademicTitleEntity, userID uint) error
	GetByID(id uint) (*entity.AcademicTitleEntity, error)
	GetByUserID(userID uint) (*entity.AcademicTitleEntity, error)
}
