package entity

type JobEntity struct {
	ID               uint   `json:"id"`
	OrganizationTin  string `json:"organization_tin"`
	OrganizationID   *uint  `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	UserID           uint   `json:"user_id"`
	PositionName     string `json:"position_name"`
}

func NewJobEntity(
	organizationID *uint,
	userID uint,
	organizationTin string,
	organizationName string,
	positionName string,
) *JobEntity {
	return &JobEntity{
		OrganizationID:   organizationID,
		UserID:           userID,
		OrganizationTin:  organizationTin,
		OrganizationName: organizationName,
		PositionName:     positionName,
	}
}

type JobWithAuthorIDEntity struct {
	*JobEntity
	AuthorID uint
}

func NewJobWithAuthorIDEntity(job *JobEntity, authorID uint) *JobWithAuthorIDEntity {
	return &JobWithAuthorIDEntity{JobEntity: job, AuthorID: authorID}
}
