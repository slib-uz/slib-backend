package response

import (
	"time"

	entity2 "slib.uz/src/core/domain/entity"
	enum2 "slib.uz/src/core/domain/entity/enum"
)

type SubmitResponse struct {
	Data []*SubmissionResponse `json:"data"`
}

func (this *SubmitResponse) ToEntityList(idempotencyKey string, applicationID uint, senderID uint, reviewers []*entity2.ReviewerEntity) []*entity2.PeerReviewSubmissionEntity {
	entities := make([]*entity2.PeerReviewSubmissionEntity, len(this.Data))

	for i, r := range reviewers {
		sub := this.findReviewerSubmission(r.ExternalID)
		entities[i] = sub.ToEntity(idempotencyKey, applicationID, r.ID, senderID)
	}
	return entities
}

func (this *SubmitResponse) findReviewerSubmission(reviewerExternalID uint) *SubmissionResponse {
	for _, s := range this.Data {
		if s.ReviewerID == reviewerExternalID {
			return s
		}
	}
	return nil
}

type SubmissionResponse struct {
	ID           uint                   `json:"id"`
	Title        string                 `json:"title"`
	SenderTitle  string                 `json:"sender_title"`
	OldDeadline  *string                `json:"old_deadline"`
	Deadline     string                 `json:"deadline"`
	Status       enum2.PeerReviewStatus `json:"status"`
	ExtraData    map[string]any         `json:"extra_data"`
	ReviewMethod enum2.PeerReviewMethod `json:"review_method"`
	ReviewerID   uint                   `json:"reviewer_id"`
}

func (this *SubmissionResponse) ToEntity(idempotencyKey string, applicationID, reviewerInternalID, senderID uint) *entity2.PeerReviewSubmissionEntity {

	var deadline time.Time
	var oldDeadline *time.Time

	if this.Deadline != "" {
		d := parseTime(this.Deadline)
		if d != nil {
			deadline = *d
		}
	}
	if this.OldDeadline != nil && *this.OldDeadline != "" {
		od := parseTime(*this.OldDeadline)
		if od != nil {
			oldDeadline = od
		}
	}

	return entity2.NewPeerReviewSubmissionEntity(
		0,
		idempotencyKey,
		this.ID,
		applicationID,
		nil,
		reviewerInternalID,
		nil,
		this.ReviewerID,
		this.SenderTitle,
		this.Title,
		this.ExtraData,
		oldDeadline,
		deadline,
		this.Status,
		this.ReviewMethod,
		senderID,
		nil,
		time.Now(),
		nil,
		0,
	)
}

func parseTime(timeStr string) *time.Time {
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil
	}
	return &parsedTime
}
