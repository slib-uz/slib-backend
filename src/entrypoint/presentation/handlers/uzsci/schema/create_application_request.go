package schema

type ApplicationAnswerRequest struct {
	FormID uint   `json:"form_id" validate:"required,min=1"`
	Value  string `json:"value"`
}

type CreateApplicationRequest struct {
	JournalID uint                       `json:"journal_id" validate:"required,min=1"`
	Answers   []ApplicationAnswerRequest `json:"answers" validate:"dive"`
}
