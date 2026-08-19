package entity

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type AntiPlagResultEntity struct {
	ID            uint               `json:"ID"`
	ReviewStageID uint               `json:"review_stage_id"`
	ReviewStage   *ReviewStageEntity `json:"review_stage"`
	ApplicationID uint               `json:"application_id"`
	Application   *ApplicationEntity `json:"application"`
	ArticleID     uint               `json:"article_id"`
	Article       *ArticleEntity     `json:"article"`
	JournalID     uint               `json:"journal_id"`
	Journal       *JournalEntity     `json:"journal"`

	ExternalID        uint                `json:"id"`
	Status            enum.AntiPlagStatus `json:"status"`
	StatusDisplay     string              `json:"status_display"`
	PlagiarismPercent float64             `json:"plagiarism_percent"`
	LegalPercent      float64             `json:"legal_percent"`
	SelfCitePercent   float64             `json:"self_cite_percent"`
	UnknownPercent    float64             `json:"unknown_percent"`
	ShortReportURL    string              `json:"short_report_url"`
	FullReportURL     string              `json:"full_report_url"`
	ExternalCreatedAt *time.Time          `json:"created_at"`
	LocalCreatedAt    *time.Time          `json:"_created_at"`

	Certificate *string `json:"certificate"`
}

func NewAntiPlagResultEntity(ID uint, reviewStageID uint, reviewStage *ReviewStageEntity, applicationID uint, application *ApplicationEntity, articleID uint, article *ArticleEntity, journalID uint, journal *JournalEntity, externalID uint, status enum.AntiPlagStatus, statusDisplay string, plagiarismPercent float64, legalPercent float64, selfCitePercent float64, unknownPercent float64, shortReportURL string, fullReportURL string, externalCreatedAt *time.Time, localCreatedAt *time.Time, certificate *string) *AntiPlagResultEntity {
	return &AntiPlagResultEntity{ID: ID, ReviewStageID: reviewStageID, ReviewStage: reviewStage, ApplicationID: applicationID, Application: application, ArticleID: articleID, Article: article, JournalID: journalID, Journal: journal, ExternalID: externalID, Status: status, StatusDisplay: statusDisplay, PlagiarismPercent: plagiarismPercent, LegalPercent: legalPercent, SelfCitePercent: selfCitePercent, UnknownPercent: unknownPercent, ShortReportURL: shortReportURL, FullReportURL: fullReportURL, ExternalCreatedAt: externalCreatedAt, LocalCreatedAt: localCreatedAt, Certificate: certificate}
}
