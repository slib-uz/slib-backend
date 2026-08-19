package entity

type PublisherEntity struct {
	ID          uint    `json:"id"`
	Tin         string  `json:"tin"`
	Name        *string `json:"name"`
	ShortName   *string `json:"short_name"`
	Logo        *string `json:"logo"`
	PhoneNumber *string `json:"phone_number"`
	FaxNumber   *string `json:"fax_number"`
	Email       *string `json:"email"`
	Website     *string `json:"website"`
	Address     *string `json:"address"`
	Description *string `json:"description"`
	IsActive    bool    `json:"is_active"`

	InstitutionID *uint              `json:"institution_id,omitempty"`
	Institution   *InstitutionEntity `json:"institution,omitempty"`

	RegionID   *uint           `json:"region_id,omitempty"`
	DistrictID *uint           `json:"district_id,omitempty"`
	Region     *RegionEntity   `json:"region,omitempty"`
	District   *DistrictEntity `json:"district,omitempty"`
}

func NewPublisherEntity(
	id uint,
	tin string,
	name *string,
	shortName *string,
	logo *string,
	phoneNumber *string,
	faxNumber *string,
	email *string,
	website *string,
	address *string,
	description *string,
	isActive bool,
) *PublisherEntity {
	return &PublisherEntity{
		ID:          id,
		Tin:         tin,
		Name:        name,
		ShortName:   shortName,
		Logo:        logo,
		PhoneNumber: phoneNumber,
		FaxNumber:   faxNumber,
		Email:       email,
		Website:     website,
		Address:     address,
		Description: description,
		IsActive:    isActive,
	}
}
