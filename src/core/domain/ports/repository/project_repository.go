package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
)

type ProjectRepository interface {
	GetAll(page, pageSize int) (*entity2.PagingEntity[entity2.ProjectEntity], error)
	GetByID(id uint) (*entity2.ProjectEntity, error)
}
