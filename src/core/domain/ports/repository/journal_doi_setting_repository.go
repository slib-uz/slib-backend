package repository

import (
	"context"

	"slib.uz/src/core/domain/entity"
)

type JournalDoiSettingRepository interface {
	GetByJournalID(ctx context.Context, journalID uint) (*entity.JournalDoiSettingEntity, error)
	Create(ctx context.Context, entity *entity.JournalDoiSettingEntity) error
	Update(ctx context.Context, entity *entity.JournalDoiSettingEntity) error
}
