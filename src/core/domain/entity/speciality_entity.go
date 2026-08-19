package entity

type SpecialityEntity struct {
	Code string            `json:"code"`
	Name map[string]string `json:"name"`
}

func NewSpecialityEntity(code string, name map[string]string) *SpecialityEntity {
	return &SpecialityEntity{Code: code, Name: name}
}
