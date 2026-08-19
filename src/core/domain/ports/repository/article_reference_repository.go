package repository

import "slib.uz/src/core/domain/entity"

type ArticleReferenceRepository interface {
	GetListByArticleID(articleID uint) ([]*entity.ReferenceEntity, error)
}
