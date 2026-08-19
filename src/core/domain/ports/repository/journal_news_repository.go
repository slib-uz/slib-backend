package repository

import "slib.uz/src/core/domain/entity"

type JournalNewsRepository interface {
	Create(news *entity.JournalNewsEntity) error
	Update(id uint, news *entity.JournalNewsEntity) error
	Delete(id uint) error
	GetByID(id uint) (*entity.JournalNewsEntity, error)
	GetByJournalID(journalID uint, page, pageSize int) (*entity.PagingEntity[entity.JournalNewsEntity], error)
}
