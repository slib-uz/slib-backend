package entity

type InstitutionEntity struct {
	ID         uint                `json:"id"`
	Name       string              `json:"name"`
	Tin        *string             `json:"tin,omitempty"`
	Logo       *string             `json:"logo,omitempty"`
	Publishers []*PublisherEntity  `json:"publishers,omitempty"`
}

func NewInstitutionEntity(id uint, name string, tin, logo *string) *InstitutionEntity {
	return &InstitutionEntity{
		ID:   id,
		Name: name,
		Tin:  tin,
		Logo: logo,
	}
}
