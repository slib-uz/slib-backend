package entity

type ArticleMessagingEntity struct {
	ArticleID uint
	Data      MessageEntity
}

func NewArticleMessagingEntity(articleID uint, data MessageEntity) *ArticleMessagingEntity {
	return &ArticleMessagingEntity{
		ArticleID: articleID,
		Data:      data,
	}
}
