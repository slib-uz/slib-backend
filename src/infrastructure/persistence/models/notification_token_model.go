package models

import (
	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity/enum"
)

type NotificationTokenModel struct {
	gorm.Model

	UserID  uint                     `gorm:"index;not null;"`
	User    *UserModel               `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Token   string                   `gorm:"size:500;not null;uniqueIndex;"`
	Segment enum.NotificationSegment `gorm:"not null;"`
}

func NewNotificationTokenModel(userID uint, token string, segment enum.NotificationSegment) *NotificationTokenModel {
	return &NotificationTokenModel{UserID: userID, Token: token, Segment: segment}
}
