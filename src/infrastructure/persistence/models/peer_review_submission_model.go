package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	enum2 "slib.uz/src/core/domain/entity/enum"
)

type PeerReviewSubmissionModel struct {
	gorm.Model

	ExternalIdempotencyKey string

	ExternalID uint `gorm:"uniqueIndex;not null"`

	ReviewerInternalID uint
	Reviewer           *ReviewerModel `gorm:"foreignKey:ReviewerInternalID;references:ID;constraint:OnDelete:SET NULL"`

	ReviewerExternalID uint

	ApplicationID uint
	Application   *ArticleApplicationModel `gorm:"foreignKey:ApplicationID;references:ID;constraint:OnDelete:SET NULL"`

	SenderTitle string `gorm:"not null;size:1024"`
	Title       string `gorm:"not null;size:4096"`

	ExtraData    datatypes.JSON
	OldDeadline  *time.Time
	Deadline     time.Time              `gorm:"not null"`
	Status       enum2.PeerReviewStatus `gorm:"not null;default:0"`
	ReviewMethod enum2.PeerReviewMethod

	SenderID uint
	Sender   *UserModel `gorm:"foreignKey:SenderID;references:ID;constraint:OnDelete:SET NULL"`

	Version int64 `gorm:"not null;default:0"`
}

func NewPeerReviewSubmissionModel(externalIdempotencyKey string, externalID uint, reviewerInternalID uint, reviewerExternalID uint, applicationID uint, senderTitle string, title string, extraData datatypes.JSON, oldDeadline *time.Time, deadline time.Time, status enum2.PeerReviewStatus, reviewMethod enum2.PeerReviewMethod, senderID uint, version int64) *PeerReviewSubmissionModel {
	return &PeerReviewSubmissionModel{ExternalIdempotencyKey: externalIdempotencyKey, ExternalID: externalID, ReviewerInternalID: reviewerInternalID, ReviewerExternalID: reviewerExternalID, ApplicationID: applicationID, SenderTitle: senderTitle, Title: title, ExtraData: extraData, OldDeadline: oldDeadline, Deadline: deadline, Status: status, ReviewMethod: reviewMethod, SenderID: senderID, Version: version}
}

func (PeerReviewSubmissionModel) TableName() string {
	return "peer_review_submissions"
}
