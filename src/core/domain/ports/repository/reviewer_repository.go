package repository

import (
	"slib.uz/src/core/domain/entity"
)

type ReviewerRepository interface {
	Create(entity *entity.ReviewerEntity) (uint, error)
	GetByScienceID(scienceID string) (*entity.ReviewerEntity, error)
	GetIDByExternalID(externalID uint) (uint, error)
	GetByID(id uint) (*entity.ReviewerEntity, error)
	GetByJournalID(journalID uint) ([]*entity.ReviewerEntity, error)
	RemoveReviewerFromJournal(journalID uint, reviewerID uint) error
	ListByIDs(ids []uint) ([]*entity.ReviewerEntity, error)
}
