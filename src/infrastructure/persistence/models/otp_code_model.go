package models

import (
	"time"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity/enum"
)

type OTPCodeModel struct {
	gorm.Model

	Phone     string          `gorm:"not null;size:50;"`
	Code      string          `gorm:"not null;size:32;"`
	SessionID string          `gorm:"not null;size:255;"`
	Purpose   enum.OTPPurpose `gorm:"size:32;not null;"`

	ExpiresAt time.Time `gorm:"not null;"`
	UsedAt    *time.Time
}

func NewOTPCodeModel(phone string, code string, sessionID string, purpose enum.OTPPurpose, expiresAt time.Time) *OTPCodeModel {
	return &OTPCodeModel{Phone: phone, Code: code, SessionID: sessionID, Purpose: purpose, ExpiresAt: expiresAt}
}

func (OTPCodeModel) TableName() string {
	return "otp_codes"
}
