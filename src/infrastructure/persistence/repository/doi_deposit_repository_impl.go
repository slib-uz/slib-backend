package repository

import (
	"context"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	_errors "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type DoiDepositRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewDoiDepositRepositoryImpl(baseRepository *BaseRepository) repository.DoiDepositRepository {
	return &DoiDepositRepositoryImpl{BaseRepository: baseRepository}
}

func (this *DoiDepositRepositoryImpl) Create(ctx context.Context, entity *entity.DoiDepositEntity) error {
	model := mapper.DoiDepositEntityToModel(entity)
	if err := this.db().WithContext(ctx).Create(model).Error; err != nil {
		return _errors.Wrap(err)
	}
	return nil
}

func (this *DoiDepositRepositoryImpl) GetByBatchID(ctx context.Context, batchID string) (*entity.DoiDepositEntity, error) {
	var _model models.DoiDepositModel
	if err := this.db().WithContext(ctx).Where("batch_id = ?", batchID).First(&_model).Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.DoiDepositModelToEntity(&_model), nil
}

func (this *DoiDepositRepositoryImpl) UpdateByBatchID(ctx context.Context, batchID string, status enum.DoiDepositStatus, message string, submissionID string, requestBody string, responseBody string) error {
	err := this.db().WithContext(ctx).
		Model(&models.DoiDepositModel{}).
		Where("batch_id = ?", batchID).
		Updates(map[string]any{
			"status":        status,
			"message":       message,
			"submission_id": submissionID,
			"request_body":  requestBody,
			"response_body": responseBody,
		}).Error
	if err != nil {
		return _errors.Wrap(err)
	}
	return nil
}

func (this *DoiDepositRepositoryImpl) List(ctx context.Context, journalIDs []uint, page, pageSize int) (*entity.PagingEntity[entity.DoiDepositEntity], error) {
	var items []*models.DoiDepositModel
	var total int64

	query := this.db().WithContext(ctx).Model(&models.DoiDepositModel{})

	if len(journalIDs) > 0 {
		query = query.
			Joins("JOIN articles ON articles.id = doi_deposits.article_id").
			Where("articles.journal_id IN ?", journalIDs)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	if err := query.
		Preload("Article").
		Order("doi_deposits.created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&items).Error; err != nil {
		return nil, _errors.Wrap(err)
	}

	return mapper.PagingMapper(items, mapper.DoiDepositModelToEntity, page, pageSize, total), nil
}
