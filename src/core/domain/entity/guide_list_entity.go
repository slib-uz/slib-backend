package entity

type GuideListEntity struct {
	ID    uint              `json:"id"`
	Title map[string]string `json:"title"`
}

func NewListGuideEntity(ID uint, title map[string]string) *GuideListEntity {
	return &GuideListEntity{ID: ID, Title: title}
}
