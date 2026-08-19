package entity

type SecretaryEntity struct {
	ID        uint   `json:"id"`
	FullName  string `json:"full_name"`
	ScienceID string `json:"science_id"`
}

func NewSecretaryEntity(ID uint, fullName string, scienceID string) *SecretaryEntity {
	return &SecretaryEntity{ID: ID, FullName: fullName, ScienceID: scienceID}
}
