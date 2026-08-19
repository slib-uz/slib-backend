package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
)

type OrganizationRepository interface {
	GetByTin(tin string) (*entity2.OrganizationEntity, error)
	GetByID(id uint) (*entity2.OrganizationEntity, error)
	GetList(page, size int, tin, name, address *string) (*entity2.PagingEntity[entity2.OrganizationEntity], error)
	Create(entity *entity2.OrganizationEntity) error
	Update(id uint, entity *entity2.OrganizationEntity) error
	Delete(id uint) error
}
