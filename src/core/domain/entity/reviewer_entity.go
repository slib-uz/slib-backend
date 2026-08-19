package entity

type ReviewerEntity struct {
	ID          uint
	ExternalID  uint             `json:"id"`
	ScienceID   string           `json:"science_id"`
	FullName    string           `json:"full_name"`
	PhoneNumber string           `json:"-"`
	Subjects    []map[string]any `json:"subjects"`
}

func NewReviewerEntity(id, externalID uint, scienceID string, fullName string, phoneNumber string, subjects []map[string]any) *ReviewerEntity {
	return &ReviewerEntity{ID: id, ExternalID: externalID, ScienceID: scienceID, FullName: fullName, PhoneNumber: phoneNumber, Subjects: subjects}
}
