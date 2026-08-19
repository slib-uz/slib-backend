package repository

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	_errors "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type JournalNewsRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewJournalNewsRepositoryImpl(baseRepository *BaseRepository) repository.JournalNewsRepository {
	return &JournalNewsRepositoryImpl{BaseRepository: baseRepository}
}

func (this *JournalNewsRepositoryImpl) Create(news *entity.JournalNewsEntity) error {
	return this.db().Create(mapper.JournalNewsEntityToModel(news)).Error
}

func (this *JournalNewsRepositoryImpl) Update(id uint, news *entity.JournalNewsEntity) error {
	model := mapper.JournalNewsEntityToModel(news)
	return this.db().Model(&models.JournalNewsModel{}).Where("id = ?", id).
		Select("title", "body", "image").
		Updates(model).Error
}

func (this *JournalNewsRepositoryImpl) Delete(id uint) error {
	return this.db().Delete(&models.JournalNewsModel{}, id).Error
}

func (this *JournalNewsRepositoryImpl) GetByID(id uint) (*entity.JournalNewsEntity, error) {
	var model models.JournalNewsModel
	if err := this.db().First(&model, id).Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.JournalNewsModelToEntity(&model), nil
}

func (this *JournalNewsRepositoryImpl) GetByJournalID(journalID uint, page, pageSize int) (*entity.PagingEntity[entity.JournalNewsEntity], error) {
	var models []*models.JournalNewsModel
	query := this.db().Model(&models).Where("journal_id = ?", journalID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&models).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	return mapper.PagingMapper(models, mapper.JournalNewsModelToEntity, page, pageSize, total), nil
}
