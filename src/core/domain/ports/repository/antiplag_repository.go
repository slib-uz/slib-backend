package repository

import (
	"time"

	entity2 "slib.uz/src/core/domain/entity"
)

type AntiPlagRepository interface {
	Create(entity *entity2.AntiPlagResultEntity) error
	GetByExternalID(externalID uint) (*entity2.AntiPlagResultEntity, error)
	GetByID(id uint) (*entity2.AntiPlagResultEntity, error)
	GetByIdWithApplication(id uint) (*entity2.AntiPlagResultEntity, error)
	Update(entity *entity2.AntiPlagResultEntity) error
	FindIDByExternalID(externalID uint) (uint, error)
	GetAllByApplicationID(applicationID uint) ([]*entity2.AntiPlagResultEntity, error)
	GetAllByReviewStageID(reviewStageID uint) ([]*entity2.AntiPlagResultEntity, error)
	GetLatestByApplicationID(applicationID uint) (*entity2.AntiPlagResultEntity, error)
	StatsByJournal(publisherID uint, startDate, endDate time.Time, page, pageSize int) (*entity2.PagingEntity[entity2.JournalAntiPlagStatsEntity], error)
	StatsByPublisher() ([]*entity2.PublisherAntiplagStatsEntity, error)
	UpdateCertificate(externalID uint, certificate *string) error
	GetByJournalID(journalID uint, page, pageSize int) (*entity2.PagingEntity[entity2.AntiPlagResultEntity], error)
}
