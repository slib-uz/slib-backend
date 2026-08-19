package entity

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type JournalApplicationEntity struct {
	ID              uint              `json:"id"`
	JournalID       *uint             `json:"journal_id"`
	Journal         *JournalEntity    `json:"journal"`
	UserID          *uint             `json:"user_id"`
	User            *UserSharedEntity `json:"user"`
	Status          enum.Status       `json:"status"`
	RejectionReason *string           `json:"rejection_reason"`
	ReviewedAt      time.Time         `json:"reviewed_at"`
}

func NewJournalApplicationEntity(ID uint, journalID *uint, journal *JournalEntity, userID *uint, user *UserSharedEntity, status enum.Status, rejectionReason *string, reviewedAt time.Time) *JournalApplicationEntity {
	return &JournalApplicationEntity{ID: ID, JournalID: journalID, Journal: journal, UserID: userID, User: user, Status: status, RejectionReason: rejectionReason, ReviewedAt: reviewedAt}
}
