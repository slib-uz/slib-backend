package entity

import (
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type NotificationEntity struct {
	ID          uint                    `json:"id"`
	UserID      *uint                   `json:"user_id"`
	Title       map[string]string       `json:"title"`
	Body        map[string]string       `json:"body"`
	Topic       *enum.NotificationTopic `json:"topic"`
	ExtraData   map[string]string       `json:"extra_data"`
	IsEmail     bool                    `json:"is_email"`
	IsSms       bool                    `json:"is_sms"`
	IsBroadcast bool                    `json:"is_broadcast"`
	IsRead      bool                    `json:"is_read"`
	CreatedAt   time.Time               `json:"created_at"`
}

func NewNotificationEntity(ID uint, userID *uint, title, body map[string]string, topic *enum.NotificationTopic, extraData map[string]string, isEmail bool, isSms bool, isBroadcast bool, isRead bool, createdAt time.Time) *NotificationEntity {
	return &NotificationEntity{ID: ID, UserID: userID, Title: title, Body: body, Topic: topic, ExtraData: extraData, IsEmail: isEmail, IsSms: isSms, IsBroadcast: isBroadcast, IsRead: isRead, CreatedAt: createdAt}
}
