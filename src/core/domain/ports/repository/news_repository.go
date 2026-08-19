package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
)

type NewsRepository interface {
	GetByID(newsID uint) (*entity2.NewsEntity, error)
	GetByPaging(page, pageSize int, ordering string, categoryID uint) (*entity2.PagingEntity[entity2.NewsEntity], error)
	Create(news *entity2.NewsEntity) error
	Update(id uint, news *entity2.NewsEntity) error
	Delete(id uint) error
	UpdateViewsCount(map[uint]int64) error
}
