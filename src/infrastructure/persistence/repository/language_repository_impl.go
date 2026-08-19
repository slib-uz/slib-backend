package repository

import (
	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type LanguageRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewLanguageRepository(baseRepository *BaseRepository) repository.LanguageRepository {
	return &LanguageRepositoryImpl{BaseRepository: baseRepository}
}

func (this *LanguageRepositoryImpl) db() *gorm.DB {
	return this.database.GormDB
}

func (this *LanguageRepositoryImpl) GetAll() ([]*entity.LanguageEntity, error) {

	var languageModels []models.LanguageModel

	if result := this.db().Find(&languageModels); result.Error != nil {
		return nil, result.Error
	}

	languages := make([]*entity.LanguageEntity, len(languageModels))
	for i, model := range languageModels {
		languages[i] = mapper.LanguageModelToEntity(&model)
	}

	return languages, nil
}

func (this *LanguageRepositoryImpl) IsExist(id uint) (bool, error) {

	var count int64
	if err := this.db().Model(&models.LanguageModel{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (this *LanguageRepositoryImpl) ExistingIds(ids []uint) ([]uint, error) {

	var languageModels []models.LanguageModel

	if result := this.db().Where("id IN ?", ids).Find(&languageModels); result.Error != nil {
		return nil, result.Error
	}

	languages := make([]uint, len(languageModels))
	for i, model := range languageModels {
		languages[i] = model.ID
	}

	return languages, nil
}
