package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
)

type JournalRatingRepository interface {
	Create(rating *entity2.JournalRatingEntity) error
	GetByJournalID(journalID uint, page, pageSize int, ordering string) (*entity2.PagingEntity[entity2.JournalRatingEntity], error)
	GetByID(id uint) (*entity2.JournalRatingEntity, error)
	GetStatsByJournalID(journalID uint) (*entity2.JournalRatingStatsEntity, error)
	DeleteByIDAndUserID(id, userID uint) error
}
