package entity

type ChatEntity struct {
	LastMessage *SupportDialogEntity `json:"last_message"`
	UnreadCount int64                `json:"unread_count"`
}

func NewChatEntity(lastMessage *SupportDialogEntity, unreadCount int64) *ChatEntity {
	return &ChatEntity{LastMessage: lastMessage, UnreadCount: unreadCount}
}
