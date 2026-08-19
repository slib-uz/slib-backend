package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
	enum2 "slib.uz/src/core/domain/entity/enum"
)

type JournalRepository interface {
	GetListByPage(page, limit int, submissionAccess enum2.AccessType, oakRegistered *bool, publisherId uint, name, description, issn, publisherName *string, languageIds, studyFieldIds []uint, fromYear, toYear *int, indexingTypes []enum2.IndexingType, sortBy, order string) (*entity2.PagingEntity[entity2.JournalBasicEntity], error)
	GetByIDWithRelations(id uint) (*entity2.JournalEntity, error)
	FindByID(id uint) (*entity2.JournalEntity, error)
	AddReviewers(journalID uint, reviewerIds []uint) error
	UpdateJournal(journalID uint, updateEntity *entity2.JournalCreateEntity) error
	GetJournalCount() (int64, error)
	GetTopJournals(page, pageSize int) (*entity2.PagingEntity[entity2.JournalEntity], error)
	UpdateViewsCount(counts map[uint]int64) error
	UpdateStatus(journalID uint, isActive bool) error
	ExistingIds(ids []uint) ([]uint, error)
	GetPublisherIdByJournalId(journalID uint) (uint, error)
	GetJournalStatistics(page, pageSize int) (*entity2.PagingEntity[entity2.JournalStatisticEntity], error)
	GetJournalStatisticsV2(page, pageSize int, institutionID, publisherID uint, name, description, issn, publisherName *string) (*entity2.PagingEntity[entity2.JournalStatisticV2Entity], error)
	GetJournalStatisticV2ByJournalID(journalID uint) (*entity2.JournalStatisticV2Entity, error)
	GetJournalsCompletionStats(journalIds []uint) ([]entity2.JournalCompletionStatsEntity, error)
}
