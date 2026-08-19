package entity

type UserSharedEntity struct {
	ID        uint   `json:"id"`
	ScienceID string `json:"science_id"`
	FullName  string `json:"full_name"`
}

func NewUserSharedEntity(ID uint, scienceID string, fullName string) *UserSharedEntity {
	return &UserSharedEntity{ID: ID, ScienceID: scienceID, FullName: fullName}
}
