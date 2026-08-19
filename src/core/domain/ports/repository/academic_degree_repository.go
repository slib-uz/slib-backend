package repository

import (
	"slib.uz/src/core/domain/entity"
)

type AcademicDegreeRepository interface {
	UpdateOrCreate(academicDegree *entity.AcademicDegreeEntity, userID uint) error
	GetByID(id uint) (*entity.AcademicDegreeEntity, error)
	GetByUserID(userID uint) (*entity.AcademicDegreeEntity, error)
}
