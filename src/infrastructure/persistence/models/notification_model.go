package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity/enum"
)

type NotificationModel struct {
	gorm.Model

	UserID *uint      `gorm:"index;"`
	User   *UserModel `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`

	Title datatypes.JSON
	Body  datatypes.JSON
	Topic *enum.NotificationTopic `gorm:"size:32;"`

	ExtraData datatypes.JSON

	IsEmail bool `gorm:"default:false;not null"`
	IsSms   bool `gorm:"default:false;not null"`

	IsBroadcast bool `gorm:"default:false;not null"`

	IsUserRead bool `gorm:"default:false;not null"` // for individual or multicast notifications
	IsRead     bool `gorm:"->; -:migration"`
}

func NewNotificationModel(userID *uint, title, body datatypes.JSON, topic *enum.NotificationTopic, extraData datatypes.JSON, isEmail bool, isSms bool, isBroadcast bool) *NotificationModel {
	return &NotificationModel{UserID: userID, Title: title, Body: body, Topic: topic, ExtraData: extraData, IsEmail: isEmail, IsSms: isSms, IsBroadcast: isBroadcast}
}
