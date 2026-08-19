package entity

type TagEntity struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Lang string `json:"lang"`
}

func NewTagEntity(ID uint, name, lang string) *TagEntity {
	return &TagEntity{ID: ID, Name: name, Lang: lang}
}
