package repository

import (
	"errors"
	"slib.uz/src/infrastructure/db"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type DegreeRepositoryImpl struct {
	database *db.Database
}

// @inject
func NewDegreeRepository(database *db.Database) repository.DegreeRepository {
	return &DegreeRepositoryImpl{
		database: database,
	}
}

func (this *DegreeRepositoryImpl) db() *gorm.DB {
	return this.database.GormDB
}

func (this *DegreeRepositoryImpl) GetByUserID(userID uint) (*entity.DegreeEntity, error) {
	var degree models.DegreeModel

	result := this.db().Where("user_id = ?", userID).First(&degree)
	if result.Error != nil {
		return nil, result.Error
	}
	return mapper.DegreeModelToEntity(&degree), nil
}

func (this *DegreeRepositoryImpl) UpdateOrCreate(userID uint, degree *entity.DegreeEntity) error {
	degreeModel := *mapper.DegreeEntityToModel(userID, degree)

	var existing *models.DegreeModel
	err := this.db().Where("user_id = ?", userID).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Not found, so Create
			if createErr := this.db().Create(&degreeModel).Error; createErr != nil {
				return createErr
			}
		} else {
			return err
		}
	} else {
		existing = mapper.DegreeUpdateMapper(existing, &degreeModel)

		if updateErr := this.db().Save(&existing).Error; updateErr != nil {
			return updateErr
		}
	}
	return nil
}
