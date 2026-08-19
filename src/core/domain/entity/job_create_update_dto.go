package entity

type JobCreateUpdateEntity struct {
	PositionName   string `json:"position_name"`
	OrganizationID *uint  `json:"organization_id"`
}

func NewJobCreateUpdateEntity(positionName string, organizationID *uint) *JobCreateUpdateEntity {
	return &JobCreateUpdateEntity{
		PositionName:   positionName,
		OrganizationID: organizationID,
	}
}
