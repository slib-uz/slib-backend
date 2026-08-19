package entity

type ReferenceEntity struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	ArticleID uint           `json:"article_id"`
	Article   *ArticleEntity `json:"article"`
}

func NewReferenceEntity(ID uint, name string, articleID uint, article *ArticleEntity) *ReferenceEntity {
	return &ReferenceEntity{ID: ID, Name: name, ArticleID: articleID, Article: article}
}
