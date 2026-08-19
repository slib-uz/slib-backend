package entity

type DistrictEntity struct {
	ID       uint              `json:"id"`
	Name     map[string]string `json:"name"`
	Soato    int               `json:"soato"`
	RegionID uint              `json:"region_id"`
}

func NewDistrictEntity(id uint, name map[string]string, soato int, regionID uint) *DistrictEntity {
	return &DistrictEntity{ID: id, Name: name, Soato: soato, RegionID: regionID}
}
