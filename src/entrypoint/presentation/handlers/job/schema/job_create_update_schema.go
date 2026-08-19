package schema

type JobCreateUpdateRequest struct {
	PositionName   string `json:"position_name"`
	OrganizationID *uint  `json:"organization_id"`
}

type JobCreateUpdateResponse struct {
	ID               uint   `json:"id"`
	OrganizationTin  string `json:"organization_tin"`
	OrganizationID   *uint  `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	UserID           uint   `json:"user_id"`
	PositionName     string `json:"position_name"`
}
