package response

type LanguageResponse struct {
	ID   uint              `json:"id"`
	Name map[string]string `json:"name"`
	Code string            `json:"code"`
}

func NewLanguageResponse(id uint, name map[string]string, code string) *LanguageResponse {
	return &LanguageResponse{ID: id, Name: name, Code: code}
}
