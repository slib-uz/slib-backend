package repository

import (
	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type NewsCategoryRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewNewsCategoryRepository(baseRepository *BaseRepository) repository.NewsCategoryRepository {
	return &NewsCategoryRepositoryImpl{BaseRepository: baseRepository}
}

func (this *NewsCategoryRepositoryImpl) db() *gorm.DB {
	return this.database.GormDB
}

func (this *NewsCategoryRepositoryImpl) GetAll() ([]*entity.NewsCategoryEntity, error) {
	var categoryModels []models.NewsCategoryModel

	if result := this.db().Find(&categoryModels); result.Error != nil {
		return nil, result.Error
	}

	categories := make([]*entity.NewsCategoryEntity, len(categoryModels))
	for i, model := range categoryModels {
		categories[i] = mapper.NewsCategoryModelToEntity(&model)
	}

	return categories, nil
}
