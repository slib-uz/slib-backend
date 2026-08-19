package mapper

import (
	"slib.uz/src/core/domain/entity"
)

func PagingMapper[T, R any](items []*T, mapper func(*T) *R, page, pageSize int, total int64) *entity.PagingEntity[R] {

	mappedItems := make([]*R, len(items))

	for i, item := range items {
		mappedItems[i] = mapper(item)
	}

	return entity.NewPagingEntity(page, pageSize, total, mappedItems)
}
