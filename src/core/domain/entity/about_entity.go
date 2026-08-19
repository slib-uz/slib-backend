package entity

type AboutEntity struct {
	ID      uint              `json:"id"`
	Content map[string]string `json:"content"`
}

func NewAboutEntity(id uint, content map[string]string) *AboutEntity {
	return &AboutEntity{ID: id, Content: content}
}
