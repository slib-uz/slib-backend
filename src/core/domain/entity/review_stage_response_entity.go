package entity

import (
	"time"

	enum2 "slib.uz/src/core/domain/entity/enum"
)

type ReviewStageResponseEntity struct {
	ID               uint                       `json:"id"`
	ApplicationID    uint                       `json:"application_id"`
	Application      *ApplicationBasicEntity    `json:"application"`
	Stage            enum2.Stage                `json:"stage" swaggertype:"integer"`
	Status           enum2.Status               `json:"status" swaggertype:"integer"`
	Reason           *string                    `json:"reason"`
	ReviewerID       *uint                      `json:"reviewer_id"`
	Reviewer         *UserBasicEntity           `json:"reviewer"`
	ReviewedAt       *time.Time                 `json:"reviewed_at"`
	CreatedAt        time.Time                  `json:"created_at"`
	PreviousID       *uint                      `json:"previous_id"`
	Previous         *ReviewStageResponseEntity `json:"previous"`
	AntiPlagResult   *AntiPlagResultEntity      `json:"anti_plag_result"`
	SpellCheckResult *SpellCheckResultEntity    `json:"spell_check_result"`
	AIDetectResult   *AiDetectResultEntity      `json:"ai_detect_result"`
	IsOld            bool                       `json:"is_old"`
	Deadline         time.Time                  `json:"deadline"`
}

func NewReviewStageResponseEntity(
	ID uint,
	applicationID uint,
	application *ApplicationBasicEntity,
	stage enum2.Stage,
	status enum2.Status,
	reason *string,
	reviewerID *uint,
	reviewer *UserBasicEntity,
	reviewedAt *time.Time,
	createdAt time.Time,
	previousID *uint,
	previous *ReviewStageResponseEntity,
	antiPlagResult *AntiPlagResultEntity,
	spellCheckResult *SpellCheckResultEntity,
	aiDetectResult *AiDetectResultEntity,
	isOld bool,
	deadline time.Time,
) *ReviewStageResponseEntity {
	return &ReviewStageResponseEntity{
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
		AntiPlagResult:   antiPlagResult,
		SpellCheckResult: spellCheckResult,
		AIDetectResult:   aiDetectResult,
		IsOld:            isOld,
		Deadline:         deadline,
	}
}
