package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
)

type PublisherRepository interface {
	GetListByPage(page, size int, name, tin string, institutionID uint, unassigned bool) (*entity2.PagingEntity[entity2.PublisherEntity], error)
	GetAll() ([]*entity2.PublisherEntity, error)
	GetByID(id uint) (*entity2.PublisherEntity, error)
	GetByTin(tin string) (*entity2.PublisherEntity, error)
	Update(id uint, entity *entity2.PublisherEntity) error
	Create(entity *entity2.PublisherEntity) error
	Delete(id uint) error
}
