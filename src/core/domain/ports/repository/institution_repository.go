package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
)

type InstitutionRepository interface {
	GetByID(id uint) (*entity2.InstitutionEntity, error)
	GetList(page, size int, tin, name *string) (*entity2.PagingEntity[entity2.InstitutionEntity], error)
	Create(entity *entity2.InstitutionEntity) error
	Update(id uint, entity *entity2.InstitutionEntity) error
	Delete(id uint) error
	SetPublishers(institutionID uint, publisherIDs []uint) error
	DetachPublishers(institutionID uint, publisherIDs []uint) (int64, error)
}
