package repository

import (
	"slib.uz/src/core/domain/entity"
)

type SecretaryRepository interface {
	GetByJournalID(journalID uint) ([]*entity.SecretaryEntity, error)
}
