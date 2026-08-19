package gateway

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type NotificationGateway interface {
	SendToTokens(id uint, title, body string, extraData map[string]string, tokens []string) (*entity.NotificationSendResultEntity, error)
	SendToTopic(title, body string, extraData map[string]string, topic enum.NotificationTopic) error
	SubscribeToTopic(topic enum.NotificationTopic, token string) error
	UnsubscribeFromTopic(topic enum.NotificationTopic, token string) error
}
