package entity

type CategoryEntity struct {
	ID   uint              `json:"id"`
	Name map[string]string `json:"name"`
}

func NewCategoryEntity(id uint, name map[string]string) *CategoryEntity {
	return &CategoryEntity{
		ID:   id,
		Name: name,
	}
}
