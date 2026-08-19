package models

import "gorm.io/gorm"

type CommentModel struct {
	gorm.Model
	ArticleID uint
	Article   *ArticleModel `gorm:"foreignKey:ArticleID;constraint:OnDelete:SET NULL;"`
	UserID    uint
	User      *UserModel `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL;"`

	Content string
	Rating  int `gorm:"type:smallint;check:rating >= 1 AND rating <= 5;not null"`
}

func NewCommentModel(articleID uint, userID uint, content string, rating int) *CommentModel {
	return &CommentModel{ArticleID: articleID, UserID: userID, Content: content, Rating: rating}
}
