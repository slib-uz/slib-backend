package schema

type InstitutionDetachPublisherRequest struct {
	PublisherIDs []uint `json:"publisher_ids" validate:"required,min=1"`
}
