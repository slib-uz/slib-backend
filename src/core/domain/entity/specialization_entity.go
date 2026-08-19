package entity

type SpecializationEntity struct {
	ID   uint              `json:"id"`
	Name map[string]string `json:"name"`
}

func NewSpecializationEntity(id uint, name map[string]string) *SpecializationEntity {
	return &SpecializationEntity{
		ID:   id,
		Name: name,
	}
}
