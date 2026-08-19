package entity

import (
	"slib.uz/src/core/domain/entity/enum"
)

type NotificationTemplateEntity struct {
	ID    uint
	Key   enum.NotificationTemplate
	Title map[string]string
	Body  map[string]string
}

func NewNotificationTemplateEntity(ID uint, key enum.NotificationTemplate, title map[string]string, body map[string]string) *NotificationTemplateEntity {
	return &NotificationTemplateEntity{ID: ID, Key: key, Title: title, Body: body}
}
