package repository

import (
	"slib.uz/src/core/domain/entity"
)

type ChiefEditorRepository interface {
	GetByJournalID(journalID uint) ([]*entity.ChiefEditorEntity, error)
}
