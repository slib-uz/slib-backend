package schema

type PublisherCreateRequest struct {
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

	RegionID   *uint `json:"region_id,omitempty"`
	DistrictID *uint `json:"district_id,omitempty"`
}
