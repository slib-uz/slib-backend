package entity

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type AuthorshipClaimEntity struct {
	ID           uint             `json:"id"`
	SenderID     uint             `json:"sender_id"`
	Sender       *UserEntity      `json:"sender"`
	ArticleID    uint             `json:"article_id"`
	Article      *ArticleEntity   `json:"article"`
	Comment      string           `json:"comment"`
	Status       enum.ClaimStatus `json:"status"`
	RejectReason string           `json:"reject_reason,omitempty"`
	ReviewedByID *uint            `json:"reviewed_by_id,omitempty"`
	ReviewedBy   *UserEntity      `json:"reviewed_by,omitempty"`
	ReviewedAt   *time.Time       `json:"reviewed_at,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

func NewAuthorshipClaimEntity(
	id uint,
	senderID uint,
	sender *UserEntity,
	articleID uint,
	article *ArticleEntity,
	comment string,
	status enum.ClaimStatus,
	createdAt time.Time,
	updatedAt time.Time,
) *AuthorshipClaimEntity {
	return &AuthorshipClaimEntity{
		ID:        id,
		SenderID:  senderID,
		Sender:    sender,
		ArticleID: articleID,
		Article:   article,
		Comment:   comment,
		Status:    status,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
