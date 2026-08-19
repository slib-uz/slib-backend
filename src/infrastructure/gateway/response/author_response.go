package response

type AuthorResponse struct {
	ExternalID uint   `json:"id"`
	FullName   string `json:"full_name"`
	ScienceID  string `json:"science_id"`
}

func NewAuthorResponse(id uint, fullName string, scienceID string) *AuthorResponse {
	return &AuthorResponse{ExternalID: id, FullName: fullName, ScienceID: scienceID}
}
