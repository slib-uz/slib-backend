package entity

type ChiefEditorEntity struct {
	ID        uint   `json:"id"`
	FullName  string `json:"full_name"`
	ScienceID string `json:"science_id"`
	Photo     string `json:"photo"`
}

func NewChiefEditorEntity(ID uint, fullName string, scienceID string, photo string) *ChiefEditorEntity {
	return &ChiefEditorEntity{ID: ID, FullName: fullName, ScienceID: scienceID, Photo: photo}
}
