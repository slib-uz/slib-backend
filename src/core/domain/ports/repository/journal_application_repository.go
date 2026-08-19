package repository

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type JournalApplicationRepository interface {
	Create(createEntity *entity.JournalCreateEntity) error
	GetListByPaging(publisherID uint, page, limit int, status *enum.Status) (*entity.PagingEntity[entity.JournalApplicationEntity], error)
	GetByISSN(issn string, issnField string) (*entity.JournalApplicationEntity, error)
	SetStatus(id uint, status enum.Status, rejectionReason *string) error
	GetByID(id uint) (*entity.JournalApplicationEntity, error)
}
