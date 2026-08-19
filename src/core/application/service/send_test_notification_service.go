package service

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
)

type SendTestNotificationUseCase struct {
	gateway gateway.NotificationGateway
}

// @inject
func NewSendTestNotificationUseCase(gateway gateway.NotificationGateway) *SendTestNotificationUseCase {
	return &SendTestNotificationUseCase{gateway: gateway}
}

func (this *SendTestNotificationUseCase) SendToTopic(topic string, title, body string) error {
	return this.gateway.SendToTopic(title, body, nil, enum.NotificationTopic(topic))
}

func (this *SendTestNotificationUseCase) SendToTokens(tokens []string, title, body string) (*entity.NotificationSendResultEntity, error) {
	return this.gateway.SendToTokens(0, title, body, nil, tokens)
}
