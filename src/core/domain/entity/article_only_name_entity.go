package entity

type ArticleOnlyNameEntity struct {
	ID   uint              `json:"id"`
	Name map[string]string `json:"name"`
}

func NewArticleOnlyNameEntity(id uint, name map[string]string) *ArticleOnlyNameEntity {
	return &ArticleOnlyNameEntity{ID: id, Name: name}
}
