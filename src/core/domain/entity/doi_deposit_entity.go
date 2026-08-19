package entity

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type DoiDepositEntity struct {
	ID           uint
	ArticleID    *uint
	BatchID      string
	DOI          string
	Status       enum.DoiDepositStatus
	Message      string
	SubmissionID string
	RequestBody  string
	ResponseBody string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewDoiDepositEntity(ID uint, articleID *uint, batchID string, doi string, status enum.DoiDepositStatus, message string, submissionID string, requestBody string, responseBody string, createdAt time.Time, updatedAt time.Time) *DoiDepositEntity {
	return &DoiDepositEntity{ID: ID, ArticleID: articleID, BatchID: batchID, DOI: doi, Status: status, Message: message, SubmissionID: submissionID, RequestBody: requestBody, ResponseBody: responseBody, CreatedAt: createdAt, UpdatedAt: updatedAt}
}
