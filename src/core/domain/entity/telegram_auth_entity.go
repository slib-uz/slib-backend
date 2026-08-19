package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"slib.uz/src/core/domain/entity/enum"
)

type TelegramAuthSessionEntity struct {
	SessionID      uuid.UUID                      `json:"session_id"`
	ScienceID      string                         `json:"science_id"`
	VerifiedPhone  string                         `json:"verified_phone"`
	TelegramChatID int64                          `json:"telegram_chat_id"`
	UserID         *uint                          `json:"user_id"`
	Status         enum.TelegramAuthSessionStatus `json:"status"`
	ExpiresAt      time.Time                      `json:"expires_at"`
}

func NewTelegramAuthSessionEntity(sessionID uuid.UUID, scienceID string, verifiedPhone string, telegramChatID int64, userID *uint, status enum.TelegramAuthSessionStatus, expiresAt time.Time) *TelegramAuthSessionEntity {
	return &TelegramAuthSessionEntity{
		SessionID:      sessionID,
		ScienceID:      scienceID,
		VerifiedPhone:  verifiedPhone,
		TelegramChatID: telegramChatID,
		UserID:         userID,
		Status:         status,
		ExpiresAt:      expiresAt,
	}
}

func NewTelegramAuthScienceIDSessionEntity(scienceID string) *TelegramAuthSessionEntity {
	return &TelegramAuthSessionEntity{
		SessionID: uuid.New(),
		ScienceID: scienceID,
		Status:    enum.TelegramAuthSessionStatusPending,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
}

func (this *TelegramAuthSessionEntity) CompleteAuth(userID uint, verifiedPhone string) error {
	if this.IsExpired() {
		return errors.New("Session is expired")
	}
	if this.Status != enum.TelegramAuthSessionStatusWaitingConfirmation && this.Status != enum.TelegramAuthSessionStatusNeedsRegistration {
		return errors.New("Invalid status for completion")
	}

	this.UserID = &userID
	this.VerifiedPhone = verifiedPhone
	this.Status = enum.TelegramAuthSessionStatusSuccess
	return nil
}

func (this *TelegramAuthSessionEntity) NeedsRegistration(verifiedPhone string) error {
	if this.IsExpired() {
		return errors.New("Session is expired")
	}
	if this.Status != enum.TelegramAuthSessionStatusWaitingConfirmation {
		return errors.New("Invalid status for registration step")
	}

	this.VerifiedPhone = verifiedPhone
	this.Status = enum.TelegramAuthSessionStatusNeedsRegistration

	this.ExpiresAt = time.Now().Add(40 * time.Minute)

	return nil
}

func (this *TelegramAuthSessionEntity) Conflict() {
	this.Status = enum.TelegramAuthSessionStatusConflict
}

func (this *TelegramAuthSessionEntity) Fail() {
	this.Status = enum.TelegramAuthSessionStatusFailed
}

func (this *TelegramAuthSessionEntity) IsExpired() bool {
	return time.Now().After(this.ExpiresAt)
}

func (this *TelegramAuthSessionEntity) WaitingForConfirmation(chatID int64) error {
	if this.IsExpired() {
		return errors.New("Session is expired")
	}
	if this.Status != enum.TelegramAuthSessionStatusPending {
		return errors.New("Invalid status for confirmation")
	}

	this.TelegramChatID = chatID
	this.Status = enum.TelegramAuthSessionStatusWaitingConfirmation
	return nil
}
