package entity

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type AiDetectResultEntity struct {
	ID            uint               `json:"ID"`
	ReviewStageID uint               `json:"review_stage_id"`
	ReviewStage   *ReviewStageEntity `json:"review_stage"`
	ApplicationID uint               `json:"application_id"`
	Application   *ApplicationEntity `json:"application"`
	ArticleID     uint               `json:"article_id"`
	Article       *ArticleEntity     `json:"article"`
	JournalID     uint               `json:"journal_id"`
	Journal       *JournalEntity     `json:"journal"`

	ExternalID    uint                `json:"id"`
	WordsCount    int                 `json:"words_count"`
	Status        enum.AntiPlagStatus `json:"status"`
	StatusDisplay string              `json:"status_display"`
	HumanPercent  float64             `json:"human_percent"`
	WarnPercent   float64             `json:"warn_percent"`
	AiPercent     float64             `json:"ai_percent"`
	ReportURL     string              `json:"report_url"`

	ExternalCreatedAt *time.Time `json:"created_at"`
	LocalCreatedAt    *time.Time `json:"_created_at"`
}

func NewAiDetectResultEntity(
	ID uint,
	reviewStageID uint,
	reviewStage *ReviewStageEntity,
	applicationID uint,
	application *ApplicationEntity,
	articleID uint,
	article *ArticleEntity,
	journalID uint,
	journal *JournalEntity,
	externalID uint,
	wordsCount int,
	status enum.AntiPlagStatus,
	statusDisplay string,
	humanPercent float64,
	warnPercent float64,
	aiPercent float64,
	reportURL string,
	externalCreatedAt *time.Time,
	localCreatedAt *time.Time,
) *AiDetectResultEntity {
	return &AiDetectResultEntity{
		ID:                ID,
		ReviewStageID:     reviewStageID,
		ReviewStage:       reviewStage,
		ApplicationID:     applicationID,
		Application:       application,
		ArticleID:         articleID,
		Article:           article,
		JournalID:         journalID,
		Journal:           journal,
		ExternalID:        externalID,
		WordsCount:        wordsCount,
		Status:            status,
		StatusDisplay:     statusDisplay,
		HumanPercent:      humanPercent,
		WarnPercent:       warnPercent,
		AiPercent:         aiPercent,
		ReportURL:         reportURL,
		ExternalCreatedAt: externalCreatedAt,
		LocalCreatedAt:    localCreatedAt,
	}
}
