package schema

type RemoveReviewerRequest struct {
	JournalID  uint `json:"journal_id"`
	ReviewerID uint `json:"reviewer_id"`
}
