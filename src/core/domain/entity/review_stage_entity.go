package entity

import (
	"time"

	enum2 "slib.uz/src/core/domain/entity/enum"
)

type ReviewStageEntity struct {
	ID               uint                    `json:"id"`
	ApplicationID    uint                    `json:"application_id"`
	Application      *ApplicationEntity      `json:"application"`
	Stage            enum2.Stage             `json:"stage" swaggertype:"integer"`
	Status           enum2.Status            `json:"status" swaggertype:"integer"`
	Reason           *string                 `json:"reason"`
	ReviewerID       *uint                   `json:"reviewer_id"`
	Reviewer         *UserEntity             `json:"reviewer"`
	ReviewedAt       *time.Time              `json:"reviewed_at"`
	CreatedAt        time.Time               `json:"created_at"`
	PreviousID       *uint                   `json:"previous_id,omitempty"`
	Previous         *ReviewStageEntity      `json:"previous,omitempty"`
	IsOld            bool                    `json:"is_old,omitempty"`
	AntiPlagResult   *AntiPlagResultEntity   `json:"anti_plag_result,omitempty"`
	SpellCheckResult *SpellCheckResultEntity `json:"spell_check_result,omitempty"`
	AIDetectResult   *AiDetectResultEntity   `json:"ai_detect_result,omitempty"`

	Deadline         time.Time  `json:"deadline"`
	ResubmitDeadline *time.Time `json:"resubmit_deadline"`
}

func NewReviewStageEntity(
	ID uint,
	applicationID uint,
	application *ApplicationEntity,
	stage enum2.Stage,
	status enum2.Status,
	reason *string,
	reviewerID *uint,
	reviewer *UserEntity,
	reviewedAt *time.Time,
	createdAt time.Time,
	previousID *uint,
	previous *ReviewStageEntity,
	isOld bool,
	antiPlagResult *AntiPlagResultEntity,
	spellCheckResult *SpellCheckResultEntity,
	aiDetectResult *AiDetectResultEntity,
	deadline time.Time,
	resubmitDeadline *time.Time,
) *ReviewStageEntity {
	return &ReviewStageEntity{
		ID:               ID,
		ApplicationID:    applicationID,
		Application:      application,
		Stage:            stage,
		Status:           status,
		Reason:           reason,
		ReviewerID:       reviewerID,
		Reviewer:         reviewer,
		ReviewedAt:       reviewedAt,
		CreatedAt:        createdAt,
		PreviousID:       previousID,
		Previous:         previous,
		IsOld:            isOld,
		AntiPlagResult:   antiPlagResult,
		SpellCheckResult: spellCheckResult,
		AIDetectResult:   aiDetectResult,
		Deadline:         deadline,
		ResubmitDeadline: resubmitDeadline,
	}
}
