package schema

type SpellCheckRequest struct {
	ApplicationID uint `json:"application_id" validate:"required"`
	ReviewStageID uint `json:"review_stage_id" validate:"required"`
}
