package schema

import (
	"slib.uz/src/core/domain/entity"
)

type CreateAuthorshipClaimRequest struct {
	ArticleIDs []uint `json:"article_ids"`
	Comment    string `json:"comment"`
}

func (this *CreateAuthorshipClaimRequest) ToEntity(userID uint) *entity.CreateAuthorshipClaimInputEntity {
	return entity.NewCreateAuthorshipClaimInputEntity(
		userID,
		this.ArticleIDs,
		this.Comment,
	)
}
