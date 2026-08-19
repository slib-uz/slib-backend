package models

import (
	"time"

	"gorm.io/gorm"
)

type SpellCheckResultModel struct {
	gorm.Model

	ReviewStageID uint
	ReviewStage   *ReviewStageModel `gorm:"foreignKey:ReviewStageID;references:ID;constraint:OnDelete:SET NULL;"`

	ApplicationID uint                     `gorm:"index"`
	Application   *ArticleApplicationModel `gorm:"constraint:OnDelete:SET NULL;"`
	JournalID     uint
	Journal       *JournalModel `gorm:"foreignKey:JournalID;references:ID;constraint:OnDelete:SET NULL;"`
	File          string        `gorm:"size:255;not null;"`
	ResultFile    *string       `gorm:"size:255;;"`
	Status        int           `gorm:"not null;default:0"` // 0 - pending, 1 - success, -1 - failed
	SubmitterID   uint          `gorm:"index"`
	Submitter     *UserModel    `gorm:"foreignKey:SubmitterID;references:ID;constraint:OnDelete:SET NULL;"`
	ResultTime    *time.Time
}

func NewSpellCheckResultModel(reviewStageID, applicationID, journalID uint, file string, resultFile *string, status int, submitterID uint) *SpellCheckResultModel {
	return &SpellCheckResultModel{ReviewStageID: reviewStageID, ApplicationID: applicationID, JournalID: journalID, File: file, ResultFile: resultFile, Status: status, SubmitterID: submitterID}
}

func (SpellCheckResultModel) TableName() string {
	return "spellcheck_results"
}
