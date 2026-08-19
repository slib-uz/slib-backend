package repository

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	_errors "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type JournalEditorialRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewJournalEditorialRepositoryImpl(baseRepository *BaseRepository) repository.JournalEditorialRepository {
	return &JournalEditorialRepositoryImpl{BaseRepository: baseRepository}
}

func (this *JournalEditorialRepositoryImpl) Create(editorial *entity.JournalEditorialEntity) error {
	return this.db().Create(mapper.JournalEditorialEntityToModel(editorial)).Error
}

func (this *JournalEditorialRepositoryImpl) Update(id uint, editorial *entity.JournalEditorialEntity) error {
	model := mapper.JournalEditorialEntityToModel(editorial)
	return this.db().Model(&models.JournalEditorialModel{}).Where("id = ?", id).
		Select("full_name", "role_code", "role_title", "photo", "science_id", "workplace", "position", "order").
		Updates(model).Error
}

func (this *JournalEditorialRepositoryImpl) Delete(id uint) error {
	return this.db().Delete(&models.JournalEditorialModel{}, id).Error
}

func (this *JournalEditorialRepositoryImpl) GetByID(id uint) (*entity.JournalEditorialEntity, error) {
	var model models.JournalEditorialModel
	if err := this.db().First(&model, id).Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.JournalEditorialModelToEntity(&model), nil
}

func (this *JournalEditorialRepositoryImpl) GetByJournalID(journalID uint, page, pageSize int) (*entity.PagingEntity[entity.JournalEditorialEntity], error) {
	var models []*models.JournalEditorialModel
	query := this.db().Model(&models).Where("journal_id = ?", journalID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	if err := query.Order(`"order" ASC`).
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&models).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	return mapper.PagingMapper(models, mapper.JournalEditorialModelToEntity, page, pageSize, total), nil
}

func (this *JournalEditorialRepositoryImpl) GetAllByJournalID(journalID uint) ([]*entity.JournalEditorialEntity, error) {
	var items []*models.JournalEditorialModel
	if err := this.db().Model(&models.JournalEditorialModel{}).
		Where("journal_id = ?", journalID).
		Order(`"order" ASC`).
		Find(&items).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	result := make([]*entity.JournalEditorialEntity, 0, len(items))
	for _, item := range items {
		result = append(result, mapper.JournalEditorialModelToEntity(item))
	}
	return result, nil
}
