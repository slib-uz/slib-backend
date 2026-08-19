package repository

import (
	"time"

	entity2 "slib.uz/src/core/domain/entity"
)

type SpellCheckResultRepository interface {
	Create(result *entity2.SpellCheckResultEntity) (*entity2.SpellCheckResultEntity, error)
	Update(id uint, status int, resultFile *string, resultTime *time.Time) error
	GetByApplicationID(applicationID uint) ([]*entity2.SpellCheckResultEntity, error)
	GetLatestByApplicationID(applicationID uint) (*entity2.SpellCheckResultEntity, error)
	StatsByJournal(publisherID uint, startDate, endDate time.Time, page, pageSize int) (*entity2.PagingEntity[entity2.JournalSpellcheckStatsEntity], error)
	StatsByPublisher() ([]*entity2.PublisherSpellcheckStatsEntity, error)
	GetByReviewStageID(reviewStageID uint) ([]*entity2.SpellCheckResultEntity, error)
	GetByIDWithApplication(id uint) (*entity2.SpellCheckResultEntity, error)
	GetByJournalID(journalID uint, page, pageSize int) (*entity2.PagingEntity[entity2.SpellCheckResultEntity], error)
}
