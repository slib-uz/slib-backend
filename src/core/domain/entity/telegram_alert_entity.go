package entity

type TelegramAlertEntity struct {
	Message string `json:"message"`
}

func NewTelegramAlertEntity(message string) *TelegramAlertEntity {
	return &TelegramAlertEntity{
		Message: message,
	}
}
