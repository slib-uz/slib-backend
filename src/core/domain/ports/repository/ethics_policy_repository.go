package repository

import (
	"slib.uz/src/core/domain/entity"
)

type EthicsPolicyRepository interface {
	GetByPaging(page, pageSize int) (*entity.PagingEntity[entity.EthicsPolicyEntity], error)
	GetByID(id uint) (*entity.EthicsPolicyEntity, error)
}
