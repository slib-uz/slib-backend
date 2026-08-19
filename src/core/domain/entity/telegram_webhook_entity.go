package entity

type TelegramWebhookUpdateEntity struct {
	UpdateID      int                          `json:"update_id"`
	Message       *TelegramMessageEntity       `json:"message"`
	CallbackQuery *TelegramCallbackQueryEntity `json:"callback_query"`
}

func NewTelegramWebhookUpdateEntity(updateID int, message *TelegramMessageEntity, callbackQuery *TelegramCallbackQueryEntity) *TelegramWebhookUpdateEntity {
	return &TelegramWebhookUpdateEntity{
		UpdateID:      updateID,
		Message:       message,
		CallbackQuery: callbackQuery,
	}
}

type TelegramMessageEntity struct {
	MessageID int                    `json:"message_id"`
	Text      string                 `json:"text"`
	Chat      TelegramChatEntity     `json:"chat"`
	Contact   *TelegramContactEntity `json:"contact"`
}

func NewTelegramMessageEntity(messageID int, text string, chat TelegramChatEntity, contact *TelegramContactEntity) *TelegramMessageEntity {
	return &TelegramMessageEntity{
		MessageID: messageID,
		Text:      text,
		Chat:      chat,
		Contact:   contact,
	}
}

type TelegramChatEntity struct {
	ID int64 `json:"id"`
}

func NewTelegramChatEntity(id int64) *TelegramChatEntity {
	return &TelegramChatEntity{
		ID: id,
	}
}

type TelegramCallbackQueryEntity struct {
	ID   string             `json:"id"`
	From TelegramChatEntity `json:"from"`
	Data string             `json:"data"`
}

func NewTelegramCallbackQueryEntity(id string, from TelegramChatEntity, data string) *TelegramCallbackQueryEntity {
	return &TelegramCallbackQueryEntity{
		ID:   id,
		From: from,
		Data: data,
	}
}

type TelegramContactEntity struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	UserID      int    `json:"user_id"`
}

func NewTelegramContactEntity(phoneNumber string, firstName string, userID int) *TelegramContactEntity {
	return &TelegramContactEntity{
		PhoneNumber: phoneNumber,
		FirstName:   firstName,
		UserID:      userID,
	}
}
