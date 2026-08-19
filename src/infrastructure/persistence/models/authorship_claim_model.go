package models

import (
	"time"

	"gorm.io/gorm"

	"slib.uz/src/core/domain/entity/enum"
)

type AuthorshipClaimModel struct {
	gorm.Model

	SenderID uint      `gorm:"not null;index:idx_one_pending_claim,unique,where:status = 'pending' AND deleted_at IS NULL"`
	Sender   UserModel `gorm:"foreignKey:SenderID"`

	ArticleID uint         `gorm:"not null;index:idx_one_pending_claim,unique,where:status = 'pending' AND deleted_at IS NULL"`
	Article   ArticleModel `gorm:"foreignKey:ArticleID"`

	Comment string `gorm:"type:text"`

	Status enum.ClaimStatus `gorm:"type:varchar(20);default:'pending'"`

	RejectReason string `gorm:"type:text"`

	ReviewedByID *uint
	ReviewedBy   *UserModel `gorm:"foreignKey:ReviewedByID"`

	ReviewedAt *time.Time
}

func (*AuthorshipClaimModel) TableName() string {
	return "authorship_claims"
}

func NewAuthorshipClaimModel(
	senderID,
	articleID uint,
	comment string,
	status enum.ClaimStatus,
) *AuthorshipClaimModel {
	return &AuthorshipClaimModel{
		SenderID:  senderID,
		ArticleID: articleID,
		Comment:   comment,
		Status:    status,
	}
}
