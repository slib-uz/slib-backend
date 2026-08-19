package entity

type ReporterEntity struct {
	ID        uint   `json:"id"`
	FullName  string `json:"full_name"`
	ScienceID string `json:"science_id"`
}

func NewReporterEntity(ID uint, fullName string, scienceID string) *ReporterEntity {
	return &ReporterEntity{ID: ID, FullName: fullName, ScienceID: scienceID}
}
