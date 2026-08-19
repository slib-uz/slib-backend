package models

import (
	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity/enum"
)

type SupportDialogModel struct {
	gorm.Model

	MessageType enum.DialogType `gorm:"not null"`
	OwnerID     uint            `gorm:"not null"`
	Owner       *UserModel      `gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;"`
	Message     string          `gorm:"not null"`
	ChatID      uint            `gorm:"not null"`
	IsRead      bool            `gorm:"not null;default:false"`
}

func NewSupportDialogModel(messageType enum.DialogType, ownerID uint, message string, chatID uint, isRead bool) *SupportDialogModel {
	return &SupportDialogModel{MessageType: messageType, OwnerID: ownerID, Message: message, ChatID: chatID, IsRead: isRead}
}

func (*SupportDialogModel) TableName() string {
	return "support_dialogs"
}
