package repository

import "slib.uz/src/core/domain/entity"

type JournalEditorialRepository interface {
	Create(editorial *entity.JournalEditorialEntity) error
	Update(id uint, editorial *entity.JournalEditorialEntity) error
	Delete(id uint) error
	GetByID(id uint) (*entity.JournalEditorialEntity, error)
	GetByJournalID(journalID uint, page, pageSize int) (*entity.PagingEntity[entity.JournalEditorialEntity], error)
	GetAllByJournalID(journalID uint) ([]*entity.JournalEditorialEntity, error)
}
