package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
)

type GuideRepository interface {
	GetByID(guideID uint) (*entity2.GuideRetrieveEntity, error)
	GetAll(page, pageSize int) (*entity2.PagingEntity[entity2.GuideListEntity], error)
}
