package logusecases

import (
	"slib.uz/src/core/domain/ports/gateway"
)

type SendTelegramAlertUseCase struct {
	telegramGateway gateway.TelegramBotGateway
}

// @inject
func NewSendTelegramAlertUseCase(telegramGateway gateway.TelegramBotGateway) *SendTelegramAlertUseCase {
	return &SendTelegramAlertUseCase{telegramGateway: telegramGateway}
}

func (uc *SendTelegramAlertUseCase) Execute(message string) error {
	return uc.telegramGateway.AlertAdmin(message)
}
