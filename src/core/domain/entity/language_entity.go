package entity

type LanguageEntity struct {
	ID   uint              `json:"id"`
	Name map[string]string `json:"name"`
	Code string            `json:"code"`
}

func NewLanguageEntity(ID uint, name map[string]string, code string) *LanguageEntity {
	return &LanguageEntity{ID: ID, Name: name, Code: code}
}
