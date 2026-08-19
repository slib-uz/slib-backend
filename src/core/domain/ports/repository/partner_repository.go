package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
)

type PartnerRepository interface {
	GetAll(page, pageSize int) (*entity2.PagingEntity[entity2.PartnerEntity], error)
}
