package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
)

type NotificationRepository interface {
	Create(notification *entity2.NotificationEntity) (uint, error)
	BulkCreate(notifications []*entity2.NotificationEntity) ([]uint, error)
	GetByID(id, userID uint) (*entity2.NotificationEntity, error)
	GetByUserIDAndBroadcast(userID uint, page, pageSize int) (*entity2.PagingEntity[entity2.NotificationEntity], error)
	GetUnreadNotificationsByUserID(userID uint, page, pageSize int) (*entity2.PagingEntity[entity2.NotificationEntity], error)
	GetUnreadCount(userID uint) (int64, error)
	Read(id, userID uint) error
}
