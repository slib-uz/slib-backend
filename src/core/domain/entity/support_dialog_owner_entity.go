package entity

type SupportDialogOwnerEntity struct {
	ID        uint   `json:"id"`
	FullName  string `json:"full_name"`
	ScienceID string `json:"science_id"`
}

func NewSupportDialogOwnerEntity(ID uint, fullName string, scienceID string) *SupportDialogOwnerEntity {
	return &SupportDialogOwnerEntity{ID: ID, FullName: fullName, ScienceID: scienceID}
}
