package enum

type TelegramAuthSessionStatus int

const (
	TelegramAuthSessionStatusConflict            TelegramAuthSessionStatus = -20
	TelegramAuthSessionStatusFailed              TelegramAuthSessionStatus = -10
	TelegramAuthSessionStatusPending             TelegramAuthSessionStatus = 10
	TelegramAuthSessionStatusWaitingConfirmation TelegramAuthSessionStatus = 20
	TelegramAuthSessionStatusSuccess             TelegramAuthSessionStatus = 30
	TelegramAuthSessionStatusNeedsRegistration   TelegramAuthSessionStatus = 40
)
