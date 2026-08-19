package models

import (
	"time"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity/enum"
)

type AntiPlagResultModel struct {
	gorm.Model

	ReviewStageID uint
	ReviewStage   *ReviewStageModel `gorm:"foreignKey:ReviewStageID;references:ID;constraint:OnDelete:SET NULL;"`

	ApplicationID uint
	Application   *ArticleApplicationModel `gorm:"foreignKey:ApplicationID;references:ID;constraint:OnDelete:SET NULL;"`

	ArticleID uint
	Article   *ArticleModel `gorm:"foreignKey:ArticleID;references:ID;constraint:OnDelete:SET NULL;"`

	JournalID uint
	Journal   *JournalModel `gorm:"foreignKey:JournalID;references:ID;constraint:OnDelete:SET NULL;"`

	ExternalID        uint                `gorm:"uniqueIndex;not null"`
	Status            enum.AntiPlagStatus `gorm:"not null"`
	StatusDisplay     string              `gorm:"size:128"`
	PlagiarismPercent float64             `gorm:"not null"`
	LegalPercent      float64             `gorm:"not null"`
	SelfCitePercent   float64             `gorm:"not null"`
	UnknownPercent    float64             `gorm:"not null"`
	ShortReportURL    string              `gorm:"size:2048"`
	FullReportURL     string              `gorm:"size:2048"`
	ExternalCreatedAt *time.Time

	Certificate *string
}

func NewAntiPlagResultModel(id, reviewStageID, applicationID, articleID, journalID, externalID uint, status enum.AntiPlagStatus, statusDisplay string, plagiarismPercent float64, legalPercent float64, selfCitePercent float64, unknownPercent float64, shortReportURL string, fullReportURL string, externalCreatedAt *time.Time) *AntiPlagResultModel {
	return &AntiPlagResultModel{Model: gorm.Model{ID: id}, ReviewStageID: reviewStageID, ApplicationID: applicationID, ArticleID: articleID, JournalID: journalID, ExternalID: externalID, Status: status, StatusDisplay: statusDisplay, PlagiarismPercent: plagiarismPercent, LegalPercent: legalPercent, SelfCitePercent: selfCitePercent, UnknownPercent: unknownPercent, ShortReportURL: shortReportURL, FullReportURL: fullReportURL, ExternalCreatedAt: externalCreatedAt}
}

func (AntiPlagResultModel) TableName() string {
	return "anti_plag_result_models"
}
