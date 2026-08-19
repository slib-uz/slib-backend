package entity

type RegionEntity struct {
	ID    uint              `json:"id"`
	Name  map[string]string `json:"name"`
	Soato int               `json:"soato"`
}

func NewRegionEntity(id uint, name map[string]string, soato int) *RegionEntity {
	return &RegionEntity{ID: id, Name: name, Soato: soato}
}
