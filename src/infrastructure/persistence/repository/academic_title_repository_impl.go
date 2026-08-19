package repository

import (
	"errors"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type AcademicTitleRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewAcademicTitleRepository(baseRepository *BaseRepository) repository.AcademicTitleRepository {
	return &AcademicTitleRepositoryImpl{BaseRepository: baseRepository}
}

func (this *AcademicTitleRepositoryImpl) UpdateOrCreate(academicTitle *entity.AcademicTitleEntity, userID uint) error {
	model := mapper.AcademicTitleEntityToModel(academicTitle)
	model.UserID = userID

	var existing models.AcademicTitleModel
	err := this.db().Where("user_id = ?", userID).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return this.db().Create(model).Error
	}
	if err != nil {
		return err
	}

	model.ID = existing.ID
	return this.db().Updates(model).Error
}

func (this *AcademicTitleRepositoryImpl) GetByID(id uint) (*entity.AcademicTitleEntity, error) {
	var academicTitle models.AcademicTitleModel
	if err := this.db().First(&academicTitle, id).Error; err != nil {
		return nil, err
	}
	return mapper.AcademicTitleModelToEntity(&academicTitle), nil
}

func (this *AcademicTitleRepositoryImpl) GetByUserID(userID uint) (*entity.AcademicTitleEntity, error) {
	var academicTitle models.AcademicTitleModel
	if err := this.db().Where("user_id = ?", userID).First(&academicTitle).Error; err != nil {
		return nil, err
	}
	return mapper.AcademicTitleModelToEntity(&academicTitle), nil
}
