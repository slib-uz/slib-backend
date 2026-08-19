package schema

type SubmitForPeerReviewRequest struct {
	ApplicationID uint   `json:"application_id"`
	ReviewerIDs   []uint `json:"reviewer_ids"`
}
