package repository

import (
	"context"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	_errors "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type JournalDoiSettingRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewJournalDoiSettingRepositoryImpl(baseRepository *BaseRepository) repository.JournalDoiSettingRepository {
	return &JournalDoiSettingRepositoryImpl{BaseRepository: baseRepository}
}

func (this *JournalDoiSettingRepositoryImpl) GetByJournalID(ctx context.Context, journalID uint) (*entity.JournalDoiSettingEntity, error) {
	var _model models.JournalDoiSettingModel
	if err := this.db().WithContext(ctx).Where("journal_id = ?", journalID).First(&_model).Error; err != nil {
		return nil, _errors.Wrap(err)
	}
	return mapper.JournalDoiSettingModelToEntity(&_model), nil
}

func (this *JournalDoiSettingRepositoryImpl) Create(ctx context.Context, e *entity.JournalDoiSettingEntity) error {
	model := mapper.JournalDoiSettingEntityToModel(e)
	if err := this.db().WithContext(ctx).Create(model).Error; err != nil {
		return _errors.Wrap(err)
	}
	return nil
}

func (this *JournalDoiSettingRepositoryImpl) Update(ctx context.Context, e *entity.JournalDoiSettingEntity) error {
	model := mapper.JournalDoiSettingEntityToModel(e)
	if err := this.db().WithContext(ctx).Model(&models.JournalDoiSettingModel{}).Where("journal_id = ?", e.JournalID).Updates(model).Error; err != nil {
		return _errors.Wrap(err)
	}
	return nil
}
