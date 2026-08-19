package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/db"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type DoctorateRepositoryImpl struct {
	database *db.Database
}

// @inject
func NewDoctorateRepository(database *db.Database) repository.DoctorateRepository {
	return &DoctorateRepositoryImpl{database: database}
}

func (this *DoctorateRepositoryImpl) db() *gorm.DB {
	return this.database.GormDB
}

func (this *DoctorateRepositoryImpl) GetByUserID(userID uint) (*entity.DoctorateEntity, error) {

	var doctorate models.DoctorateModel
	result := this.db().Where("user_id = ? and status_code = ?", userID, 500).First(&doctorate)

	if result.Error != nil {
		return nil, infraError.Wrap(result.Error)
	}
	return mapper.DoctorateModelToEntity(&doctorate), nil
}

func (this *DoctorateRepositoryImpl) BulkCreate(userID uint, doctorates []*entity.DoctorateEntity) error {
	var docModels []models.DoctorateModel
	for _, item := range doctorates {
		model := mapper.DoctorateEntityToModel(userID, item)
		docModels = append(docModels, *model)
	}

	return this.db().Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "external_id"}},
			DoUpdates: clause.AssignmentColumns(mapper.DoctorateModelFields()),
		},
	).Create(&docModels).Error
}
