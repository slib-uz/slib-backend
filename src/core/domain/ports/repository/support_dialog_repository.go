package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
)

type SupportDialogRepository interface {
	CreateQuestion(supportDialog *entity2.SupportDialogEntity) error
	CreateAnswer(supportDialog *entity2.SupportDialogEntity) error
	GetByPaging(page, pageSize int, ordering string, userID uint) (*entity2.PagingEntity[entity2.SupportDialogEntity], error)
	GetByChatID(page, pageSize int, ordering string, chatID uint) (*entity2.PagingEntity[entity2.SupportDialogEntity], error)
	GetChatsByPaging(page, pageSize int, ordering string) (*entity2.PagingEntity[entity2.ChatEntity], error)
	GetUnreadCount(chatID uint) (int64, error)
	MarkAsRead(chatID uint) error
}
