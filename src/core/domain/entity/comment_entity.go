package entity

import "time"

type CommentEntity struct {
	ID        uint                `json:"id"`
	ArticleID uint                `json:"article_id"`
	Article   *ArticleInputEntity `json:"article"`
	UserID    uint                `json:"user_id"`
	User      *UserSharedEntity   `json:"user"`
	Content   string              `json:"content"`
	Rating    int                 `json:"rating"`
	CreatedAt time.Time           `json:"created_at"`
}

func NewCommentEntity(ID uint, articleID uint, article *ArticleInputEntity, userID uint, user *UserSharedEntity, content string, rating int, createdAt time.Time) *CommentEntity {
	return &CommentEntity{ID: ID, ArticleID: articleID, Article: article, UserID: userID, User: user, Content: content, Rating: rating, CreatedAt: createdAt}
}
