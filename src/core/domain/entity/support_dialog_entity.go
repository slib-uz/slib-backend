package entity

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type SupportDialogEntity struct {
	ID          uint                      `json:"id"`
	MessageType enum.DialogType           `json:"message_type"`
	OwnerID     uint                      `json:"owner_id"`
	Owner       *SupportDialogOwnerEntity `json:"owner"`
	Message     string                    `json:"message"`
	ChatID      uint                      `json:"chat_id"`
	IsRead      bool                      `json:"is_read"`
	CreatedAt   time.Time                 `json:"created_at"`
}

func NewSupportDialogEntity(ID uint, messageType enum.DialogType, ownerID uint, owner *SupportDialogOwnerEntity, message string, chatID uint, isRead bool, createdAt time.Time) *SupportDialogEntity {
	return &SupportDialogEntity{ID: ID, MessageType: messageType, OwnerID: ownerID, Owner: owner, Message: message, ChatID: chatID, IsRead: isRead, CreatedAt: createdAt}
}
