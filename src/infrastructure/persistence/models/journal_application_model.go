package models

import (
	"time"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity/enum"
)

type JournalApplicationModel struct {
	gorm.Model

	JournalID *uint         `gorm:"index"`
	Journal   *JournalModel `gorm:"foreignKey:JournalID;constraint:OnDelete:SET NULL;"`

	UserID *uint      `gorm:"index"`
	User   *UserModel `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL;"`

	Status          enum.Status `gorm:"type:integer;default:0"`
	RejectionReason *string     `gorm:"size:2048"`
	ReviewedAt      time.Time
}

func (JournalApplicationModel) TableName() string {
	return "journal_applications"
}
