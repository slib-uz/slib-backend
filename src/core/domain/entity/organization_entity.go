package entity

type OrganizationEntity struct {
	ID      uint                   `json:"id"`
	Name    string                 `json:"name"`
	Tin     *string                `json:"tin"`
	Address *string                `json:"address,omitempty"`
	SoatoID *uint                  `json:"soato_id,omitempty"`
	Soato   *SoatoClassifierEntity `json:"soato,omitempty"`
}

func NewOrganizationEntity(id uint, name string, tin *string, address *string, soatoID *uint) *OrganizationEntity {
	return &OrganizationEntity{
		ID:      id,
		Name:    name,
		Tin:     tin,
		Address: address,
		SoatoID: soatoID,
	}
}
