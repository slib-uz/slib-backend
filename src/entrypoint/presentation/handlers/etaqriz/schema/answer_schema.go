package schema

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type PeerReviewAnswerSchema struct {
	Version int64 `json:"version"`

	SubmissionID uint                  `json:"submission_id"`
	Status       enum.PeerReviewStatus `json:"status"`
	OldDeadline  *time.Time            `json:"old_deadline"`
	Deadline     time.Time             `json:"deadline"`
	ReviewerID   uint                  `json:"reviewer_id"`
	CreatedAt    time.Time             `json:"created_at"`
}
