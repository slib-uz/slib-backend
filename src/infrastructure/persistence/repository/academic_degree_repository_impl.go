package repository

import (
	"errors"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type AcademicDegreeRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewAcademicDegreeRepository(baseRepository *BaseRepository) repository.AcademicDegreeRepository {
	return &AcademicDegreeRepositoryImpl{BaseRepository: baseRepository}
}

func (r *AcademicDegreeRepositoryImpl) UpdateOrCreate(academicDegree *entity.AcademicDegreeEntity, userID uint) error {
	model := mapper.AcademicDegreeEntityToModel(academicDegree)
	model.UserID = userID

	var existing models.AcademicDegreeModel
	err := r.db().Where("user_id = ?", userID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db().Create(model).Error
	}
	if err != nil {
		return err
	}

	model.ID = existing.ID
	return r.db().Updates(model).Error
}

func (this *AcademicDegreeRepositoryImpl) GetByID(id uint) (*entity.AcademicDegreeEntity, error) {
	var academicDegree models.AcademicDegreeModel
	if err := this.db().First(&academicDegree, id).Error; err != nil {
		return nil, err
	}
	return mapper.AcademicDegreeModelToEntity(&academicDegree), nil
}

func (this *AcademicDegreeRepositoryImpl) GetByUserID(userID uint) (*entity.AcademicDegreeEntity, error) {
	var academicDegree models.AcademicDegreeModel
	if err := this.db().Where("user_id = ?", userID).First(&academicDegree).Error; err != nil {
		return nil, err
	}
	return mapper.AcademicDegreeModelToEntity(&academicDegree), nil
}
