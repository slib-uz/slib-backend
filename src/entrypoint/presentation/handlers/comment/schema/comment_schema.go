package schema

import (
	"time"

	"slib.uz/src/core/domain/entity"
)

type CommentCreateRequest struct {
	ArticleID uint   `json:"article_id" validate:"required"`
	Content   string `json:"content"`
	Rating    int    `json:"rating" validate:"required,gt=0,lt=6"` // Rating must be between 1 and 5
}

func (this *CommentCreateRequest) ToEntity() *entity.CommentEntity {
	return entity.NewCommentEntity(
		0,
		this.ArticleID,
		nil,
		0,
		nil,
		this.Content,
		this.Rating,
		time.Now(),
	)
}
