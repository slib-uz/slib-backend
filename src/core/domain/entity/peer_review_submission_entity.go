package entity

import (
	"time"

	enum2 "slib.uz/src/core/domain/entity/enum"
)

type PeerReviewSubmissionEntity struct {
	ID uint `json:"id"`

	ExternalIdempotencyKey string `json:"external_idempotency_key"`

	ExternalID uint `json:"external_id"`

	ApplicationID uint               `json:"application_id"`
	Application   *ApplicationEntity `json:"application"`

	ReviewerInternalID uint            `json:"reviewer_internal_id"`
	Reviewer           *ReviewerEntity `json:"reviewer"`

	ReviewerExternalID uint `json:"reviewer_external_id"`

	SenderTitle string `json:"sender_title"`
	Title       string `json:"title"`

	ExtraData    map[string]any         `json:"extra_data"`
	OldDeadline  *time.Time             `json:"old_deadline"`
	Deadline     time.Time              `json:"deadline"`
	Status       enum2.PeerReviewStatus `json:"status" swaggertype:"integer"`
	ReviewMethod enum2.PeerReviewMethod `json:"review_method" swaggertype:"integer"`

	SenderID uint        `json:"sender_id"`
	Sender   *UserEntity `json:"sender"`

	CreatedAt time.Time `json:"created_at"`

	Reviews []*PeerReviewResultEntity `json:"reviews"`

	Version int64 `json:"version"`
}

func NewPeerReviewSubmissionEntity(ID uint, externalIdempotencyKey string, externalID uint, applicationID uint, application *ApplicationEntity, reviewerInternalID uint, reviewer *ReviewerEntity, reviewerExternalID uint, senderTitle string, title string, extraData map[string]any, oldDeadline *time.Time, deadline time.Time, status enum2.PeerReviewStatus, reviewMethod enum2.PeerReviewMethod, senderID uint, sender *UserEntity, createdAt time.Time, reviews []*PeerReviewResultEntity, version int64) *PeerReviewSubmissionEntity {
	return &PeerReviewSubmissionEntity{ID: ID, ExternalIdempotencyKey: externalIdempotencyKey, ExternalID: externalID, ApplicationID: applicationID, Application: application, ReviewerInternalID: reviewerInternalID, Reviewer: reviewer, ReviewerExternalID: reviewerExternalID, SenderTitle: senderTitle, Title: title, ExtraData: extraData, OldDeadline: oldDeadline, Deadline: deadline, Status: status, ReviewMethod: reviewMethod, SenderID: senderID, Sender: sender, CreatedAt: createdAt, Reviews: reviews, Version: version}
}
