package schema

import (
	"slib.uz/src/core/domain/entity"
)

type ArticleAuthorAffiliationCreateRequest struct {
	AuthorID       uint   `json:"author_id"`
	OrganizationID *uint  `json:"organization_id"`
	PositionName   string `json:"position_name"`
}

func (this *ArticleAuthorAffiliationCreateRequest) ToEntity() *entity.ArticleAuthorAffiliationEntity {
	return entity.NewArticleAuthorAffiliationEntity(
		0,
		nil,
		this.AuthorID,
		this.OrganizationID,
		"",
		"",
		this.PositionName,
	)
}
