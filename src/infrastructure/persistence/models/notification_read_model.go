package models

import "gorm.io/gorm"

// only broadcast notifications

type NotificationReadModel struct {
	gorm.Model

	NotificationID uint               `gorm:"not null;uniqueIndex:idx_notification_id_user_id"`
	Notification   *NotificationModel `gorm:"foreignKey:NotificationID;constraint:OnDelete:CASCADE;"`
	UserID         uint               `gorm:"not null;uniqueIndex:idx_notification_id_user_id"`
	User           *UserModel         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

func NewNotificationReadModel(notificationID, userID uint) *NotificationReadModel {
	return &NotificationReadModel{NotificationID: notificationID, UserID: userID}
}
