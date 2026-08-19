package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type NotificationSendResultModel struct {
	gorm.Model
	NotificationID uint `gorm:"not null;index"`
	SuccessCount   int  `gorm:"not null;default:0"`
	FailureCount   int  `gorm:"not null;default:0"`
	FailedTokens   datatypes.JSON
	Errors         datatypes.JSON
}

func NewNotificationSendResultModel(notificationID uint, successCount int, failureCount int, failedTokens datatypes.JSON, errors datatypes.JSON) *NotificationSendResultModel {
	return &NotificationSendResultModel{NotificationID: notificationID, SuccessCount: successCount, FailureCount: failureCount, FailedTokens: failedTokens, Errors: errors}
}
