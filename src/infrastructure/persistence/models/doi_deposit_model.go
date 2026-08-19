package models

import (
	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity/enum"
)

type DoiDepositModel struct {
	gorm.Model
	ArticleID    *uint                 `gorm:"index;"`
	Article      *ArticleModel         `gorm:"constraint:OnDelete:SET NULL;"`
	BatchID      string                `gorm:"uniqueIndex;not null"`
	DOI          string                `gorm:"not null"`
	Status       enum.DoiDepositStatus `gorm:"not null;default:PENDING"`
	Message      string                `gorm:"type:text"`
	SubmissionID string
	RequestBody  string `gorm:"type:text"`
	ResponseBody string `gorm:"type:text"`
}

func NewDoiDepositModel(articleID *uint, batchID string, doi string, status enum.DoiDepositStatus, message string, submissionID string, requestBody string, responseBody string) *DoiDepositModel {
	return &DoiDepositModel{ArticleID: articleID, BatchID: batchID, DOI: doi, Status: status, Message: message, SubmissionID: submissionID, RequestBody: requestBody, ResponseBody: responseBody}
}

func (*DoiDepositModel) TableName() string {
	return "doi_deposits"
}
