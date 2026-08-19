package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
)

type StudyFieldRepository interface {
	GetAll(search string) ([]*entity2.StudyFieldEntity, error)
	GetByPaging(page, pageSize int, search string) (*entity2.PagingEntity[entity2.StudyFieldEntity], error)
	GetByJournalID(journalID uint) ([]*entity2.StudyFieldEntity, error)
	ExistingIds(ids []uint) ([]uint, error)
	Create(*entity2.StudyFieldEntity) error
	Delete(id uint) error
	Update(*entity2.StudyFieldEntity) error
}
