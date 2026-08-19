package repository

import (
	"time"

	entity2 "slib.uz/src/core/domain/entity"
)

type ArticlePurchaseRepository interface {
	Create(purchase *entity2.ArticlePurchaseEntity) error
	IsExists(articleID, userID uint) (bool, error)
	GetByJournalID(journalID uint, page, pageSize int, startDate, endDate time.Time) (*entity2.PagingEntity[entity2.ArticlePurchaseEntity], error)
	StatsByJournal(publisherID uint, startDate, endDate time.Time, page, pageSize int) (*entity2.PagingEntity[entity2.JournalArticlePurchaseStatsEntity], error)
	GetByUserID(userID uint, page, pageSize int, search string) (*entity2.PagingEntity[entity2.ArticlePurchaseEntity], error)
}
