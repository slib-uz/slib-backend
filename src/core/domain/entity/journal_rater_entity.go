package entity

type JournalRaterEntity struct {
	ID        uint   `json:"id"`
	FullName  string `json:"full_name"`
	ScienceID string `json:"science_id"`
}

func NewJournalRaterEntity(ID uint, fullName string, scienceID string) *JournalRaterEntity {
	return &JournalRaterEntity{ID: ID, FullName: fullName, ScienceID: scienceID}
}
