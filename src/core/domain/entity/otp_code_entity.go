package entity

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type OTPCodeEntity struct {
	ID uint

	Phone     string
	Code      string
	SessionID string
	Purpose   enum.OTPPurpose

	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

func NewOTPCodeEntity(ID uint, phone string, code string, sessionID string, purpose enum.OTPPurpose, expiresAt time.Time, usedAt *time.Time, createdAt time.Time) *OTPCodeEntity {
	return &OTPCodeEntity{ID: ID, Phone: phone, Code: code, SessionID: sessionID, Purpose: purpose, ExpiresAt: expiresAt, UsedAt: usedAt, CreatedAt: createdAt}
}
