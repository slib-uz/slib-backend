package entity

type NewsCategoryEntity struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewNewsCategoryEntity(id uint, name string) *NewsCategoryEntity {
	return &NewsCategoryEntity{
		ID:   id,
		Name: name,
	}
}
