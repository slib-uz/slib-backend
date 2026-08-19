package repository

import (
	"context"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	_errors "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type EditionRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewEditionRepositoryImpl(baseRepository *BaseRepository) repository.EditionRepository {
	return &EditionRepositoryImpl{BaseRepository: baseRepository}
}

func (this *EditionRepositoryImpl) Create(edition *entity.EditionEntity) error {
	_model := mapper.EditionEntityToModel(edition)
	return this.db().Create(_model).Error
}

func (this *EditionRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.EditionEntity, error) {
	var _model models.EditionModel
	if err := this.db().WithContext(ctx).Preload("Journal").First(&_model, id).Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.EditionModelToEntity(&_model), nil
}

func (this *EditionRepositoryImpl) Update(ctx context.Context, edition *entity.EditionEntity) error {
	return this.db().WithContext(ctx).
		Model(&models.EditionModel{}).
		Where("id = ?", edition.ID).
		Updates(map[string]interface{}{
			"file":         edition.File,
			"cover_image":  edition.CoverImage,
			"name":         edition.Name,
			"number":       edition.Number,
			"volume":       edition.Volume,
			"published_at": edition.PublishedAt,
		}).Error
}

func (this *EditionRepositoryImpl) AttachArticles(ctx context.Context, editionID uint, articleIDs []uint) (int64, error) {
	result := this.db().WithContext(ctx).
		Model(&models.ArticleModel{}).
		Where("id IN ? AND edition_id IS NULL", articleIDs).
		Update("edition_id", editionID)

	if result.Error != nil {
		return 0, _errors.Wrap(result.Error)
	}

	return result.RowsAffected, nil
}

func (this *EditionRepositoryImpl) DetachArticles(ctx context.Context, editionID uint, articleIDs []uint) (int64, error) {
	result := this.db().WithContext(ctx).
		Model(&models.ArticleModel{}).
		Where("id IN ? AND edition_id = ?", articleIDs, editionID).
		Update("edition_id", nil)

	if result.Error != nil {
		return 0, _errors.Wrap(result.Error)
	}

	return result.RowsAffected, nil
}

func (this *EditionRepositoryImpl) Delete(ctx context.Context, editionID uint) error {
	return this.db().WithContext(ctx).Delete(&models.EditionModel{}, editionID).Error
}

func (this *EditionRepositoryImpl) GetByJournalID(ctx context.Context, journalID uint, page int, pageSize int, search string, year int) (*entity.PagingEntity[entity.EditionEntity], error) {
	var _models []*models.EditionModel

	query := this.db().WithContext(ctx).Model(&models.EditionModel{})

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if year > 0 {
		query = query.Where("EXTRACT(YEAR FROM published_at) = ?", year)
	}

	var total int64
	if err := query.Where("journal_id = ?", journalID).Count(&total).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&_models).Error
	if err != nil {
		return nil, _errors.Wrap(err)
	}

	return mapper.PagingMapper(_models, mapper.EditionModelToEntity, page, pageSize, total), nil
}
