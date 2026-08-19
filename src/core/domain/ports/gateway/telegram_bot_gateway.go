package gateway

type TelegramBotGateway interface {
	AlertAdmin(message string) error
}
