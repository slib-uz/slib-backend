package repository

import (
	"slib.uz/src/core/domain/entity"
)

type ArticleAuthorAffiliationRepository interface {
	Create(articleAuthorAffiliation *entity.ArticleAuthorAffiliationEntity) (*entity.ArticleAuthorAffiliationEntity, error)
	GetByArticleID(articleID uint) ([]*entity.ArticleAuthorAffiliationEntity, error)
	ExistingIds(ids []uint) ([]uint, error)
	BulkUpdateArticleIds(articleID uint, articleAuthorAffiliationIDs []uint) ([]uint, error)
	GetAuthorIdsByArticleAuthorAffiliationIDs(articleAuthorAffiliationIDs []uint) ([]uint, error)
	BulkCreateForAuthorsNeedArticleAuthorAffiliation(affiliations []*entity.ArticleAuthorAffiliationEntity) error
	GetListByAuthorIDs(authorIDs []uint) ([]*entity.ArticleAuthorAffiliationEntity, error)
}
