package repository

import "slib.uz/src/core/domain/entity"

type AuthorshipClaimRepository interface {
	Create(claim *entity.AuthorshipClaimEntity) error
	CreateBatch(claims []*entity.AuthorshipClaimEntity) error
	GetList(page, size int, filters map[string]interface{}) (*entity.PagingEntity[entity.AuthorshipClaimEntity], error)
	GetByID(id uint) (*entity.AuthorshipClaimEntity, error)
	FindPendingByArticleAndUser(articleID, userID uint) (*entity.AuthorshipClaimEntity, error)
	FindPendingByArticleIDsAndUser(articleIDs []uint, userID uint) ([]*entity.AuthorshipClaimEntity, error)
	Update(claim *entity.AuthorshipClaimEntity) error
}
