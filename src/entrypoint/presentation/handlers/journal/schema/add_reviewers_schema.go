package schema

type AddReviewersRequest struct {
	JournalID   uint   `json:"journal_id"`
	ReviewerIds []uint `json:"reviewer_ids"`
}
