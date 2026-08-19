package models

import (
	"time"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity/enum"
)

type AiDetectResultModel struct {
	gorm.Model

	ReviewStageID uint
	ReviewStage   *ReviewStageModel `gorm:"foreignKey:ReviewStageID;references:ID;constraint:OnDelete:SET NULL;"`

	ApplicationID uint
	Application   *ArticleApplicationModel `gorm:"foreignKey:ApplicationID;references:ID;constraint:OnDelete:SET NULL;"`

	ArticleID uint
	Article   *ArticleModel `gorm:"foreignKey:ArticleID;references:ID;constraint:OnDelete:SET NULL;"`

	JournalID uint
	Journal   *JournalModel `gorm:"foreignKey:JournalID;references:ID;constraint:OnDelete:SET NULL;"`

	ExternalID    uint                `gorm:"uniqueIndex;not null"`
	WordsCount    int                 `gorm:"not null;default:0"`
	Status        enum.AntiPlagStatus `gorm:"not null"`
	StatusDisplay string              `gorm:"size:128"`
	HumanPercent  float64             `gorm:"not null;default:0"`
	WarnPercent   float64             `gorm:"not null;default:0"`
	AiPercent     float64             `gorm:"not null;default:0"`
	ReportURL     string              `gorm:"size:2048"`

	ExternalCreatedAt *time.Time
}

func NewAiDetectResultModel(id, reviewStageID, applicationID, articleID, journalID, externalID uint, wordsCount int, status enum.AntiPlagStatus, statusDisplay string, humanPercent, warnPercent, aiPercent float64, reportURL string, externalCreatedAt *time.Time) *AiDetectResultModel {
	return &AiDetectResultModel{
		Model:             gorm.Model{ID: id},
		ReviewStageID:     reviewStageID,
		ApplicationID:     applicationID,
		ArticleID:         articleID,
		JournalID:         journalID,
		ExternalID:        externalID,
		WordsCount:        wordsCount,
		Status:            status,
		StatusDisplay:     statusDisplay,
		HumanPercent:      humanPercent,
		WarnPercent:       warnPercent,
		AiPercent:         aiPercent,
		ReportURL:         reportURL,
		ExternalCreatedAt: externalCreatedAt,
	}
}

func (AiDetectResultModel) TableName() string {
	return "ai_detect_result_models"
}
