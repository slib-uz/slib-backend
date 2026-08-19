package repository

import (
	"slib.uz/src/core/domain/entity"
)

type NotificationTokenRepository interface {
	CreateOrUpdate(entity *entity.NotificationTokenEntity) (isCreated bool, err error)
	Delete(userID uint, token string) error
	GetTokensByUserID(userID uint) ([]string, error)
}
