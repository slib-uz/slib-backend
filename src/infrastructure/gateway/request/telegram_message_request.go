package request

type TelegramKeyboardButtonRequest struct {
	Text           string `json:"text"`
	RequestContact bool   `json:"request_contact,omitempty"`
}

func NewTelegramKeyboardButtonRequest(text string, requestContact bool) TelegramKeyboardButtonRequest {
	return TelegramKeyboardButtonRequest{Text: text, RequestContact: requestContact}
}

type TelegramReplyKeyboardMarkupRequest struct {
	Keyboard        [][]TelegramKeyboardButtonRequest `json:"keyboard"`
	ResizeKeyboard  bool                              `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard bool                              `json:"one_time_keyboard,omitempty"`
}

func NewTelegramReplyKeyboardMarkupRequest(keyboard [][]TelegramKeyboardButtonRequest) TelegramReplyKeyboardMarkupRequest {
	return TelegramReplyKeyboardMarkupRequest{
		Keyboard:        keyboard,
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}
}

type TelegramSendContactRequest struct {
	ChatID      int64                              `json:"chat_id"`
	Text        string                             `json:"text"`
	ReplyMarkup TelegramReplyKeyboardMarkupRequest `json:"reply_markup"`
}

func NewTelegramSendContactRequest(chatID int64, text string, markup TelegramReplyKeyboardMarkupRequest) TelegramSendContactRequest {
	return TelegramSendContactRequest{ChatID: chatID, Text: text, ReplyMarkup: markup}
}

type TelegramSendMessageRequest struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func NewTelegramSendMessageRequest(chatID int64, text string) TelegramSendMessageRequest {
	return TelegramSendMessageRequest{ChatID: chatID, Text: text}
}
