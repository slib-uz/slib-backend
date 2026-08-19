package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
)

type CommentRepository interface {
	Create(comment *entity2.CommentEntity) error
	GetByArticleID(articleID uint, page, pageSize int) (*entity2.PagingEntity[entity2.CommentEntity], error)
	GetByUserID(userID uint, page, pageSize int) (*entity2.PagingEntity[entity2.CommentEntity], error)
	GetStatsByArticleID(articleID uint) (*entity2.CommentStatsEntity, error)
	DeleteByIDAndUserID(id, userID uint) error
}
