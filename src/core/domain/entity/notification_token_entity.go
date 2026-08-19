package entity

import (
	"slib.uz/src/core/domain/entity/enum"
)

type NotificationTokenEntity struct {
	ID      uint                     `json:"id"`
	UserID  uint                     `json:"user_id"`
	Token   string                   `json:"token"`
	Segment enum.NotificationSegment `json:"segment"`
}

func NewNotificationTokenEntity(ID uint, userID uint, token string, segment enum.NotificationSegment) *NotificationTokenEntity {
	return &NotificationTokenEntity{ID: ID, UserID: userID, Token: token, Segment: segment}
}
