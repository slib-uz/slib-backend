package entity

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type PeerReviewResultEntity struct {
	ExternalID           uint `json:"id"`
	SubmissionExternalID uint `json:"submission_id"`

	Status             enum.PeerReviewStatus `json:"status"`
	Review             map[string]string     `json:"review"`
	Comment            string                `json:"comment"`
	ReportFileLink     string                `json:"report_file_link"`
	SubmissionFileLink string                `json:"submission_file_link"`

	ExternalCreatedAt time.Time `json:"created_at"`
}

func NewPeerReviewResultEntity(externalID uint, submissionExternalID uint, status enum.PeerReviewStatus, review map[string]string, comment string, reportFileLink string, submissionFileLink string, externalCreatedAt time.Time) *PeerReviewResultEntity {
	return &PeerReviewResultEntity{ExternalID: externalID, SubmissionExternalID: submissionExternalID, Status: status, Review: review, Comment: comment, ReportFileLink: reportFileLink, SubmissionFileLink: submissionFileLink, ExternalCreatedAt: externalCreatedAt}
}
