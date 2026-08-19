package schema

type ApplicationPublishSchema struct {
	ApplicationID uint   `json:"application_id"`
	FinalFile     string `json:"final_file" validate:"required"`
}
