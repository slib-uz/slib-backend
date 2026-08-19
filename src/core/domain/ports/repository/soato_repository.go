package repository

import "slib.uz/src/core/domain/entity"

type SoatoRepository interface {
	GetAll(page, size int, isChildren bool) (*entity.PagingEntity[entity.SoatoClassifierEntity], error)
}
