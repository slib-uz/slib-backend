package repository

import (
	"time"

	entity2 "slib.uz/src/core/domain/entity"
)

type ArticleSubmissionRepository interface {
	Create(article *entity2.ArticleCreateEntity,
		authorAffiliations []*entity2.ArticleAuthorAffiliationEntity,
		updateAuthorAffiliationIDs []uint,
		userId uint,
		technicalReviewDeadline time.Time) (*entity2.ApplicationEntity, error)
}
