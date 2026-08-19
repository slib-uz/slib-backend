package models

import (
	"time"

	"gorm.io/gorm"
	enum2 "slib.uz/src/core/domain/entity/enum"
)

type ReviewStageModel struct {
	gorm.Model

	ApplicationID uint                     `gorm:"index"`
	Application   *ArticleApplicationModel `gorm:"foreignKey:ApplicationID;references:ID;constraint:OnDelete:SET NULL;"`

	Stage  enum2.Stage  `gorm:"size:32;not null;"`
	Status enum2.Status `gorm:"not null;default:0;"`
	Reason *string

	ReviewerID *uint
	Reviewer   *UserModel `gorm:"foreignKey:ReviewerID;constraint:OnDelete:SET NULL;"`
	ReviewedAt *time.Time

	Deadline         time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	ResubmitDeadline *time.Time

	PreviousID *uint
	Previous   *ReviewStageModel `gorm:"foreignKey:PreviousID;references:ID;constraint:OnDelete:SET NULL;"`
	IsOld      bool              `gorm:"default:false;not null;"`

	AntiPlagResults   []*AntiPlagResultModel   `gorm:"foreignKey:ReviewStageID;"`
	SpellCheckResults []*SpellCheckResultModel `gorm:"foreignKey:ReviewStageID;"`
	AiDetectResults   []*AiDetectResultModel   `gorm:"foreignKey:ReviewStageID;"`
}

func NewReviewStageModel(id, applicationID uint, stage enum2.Stage, status enum2.Status, previousID *uint, deadline time.Time, resubmitDeadline *time.Time) ReviewStageModel {
	return ReviewStageModel{Model: gorm.Model{ID: id}, ApplicationID: applicationID, Stage: stage, Status: status, Deadline: deadline, ResubmitDeadline: resubmitDeadline, PreviousID: previousID}
}

func (ReviewStageModel) TableName() string {
	return "review_stages"
}
