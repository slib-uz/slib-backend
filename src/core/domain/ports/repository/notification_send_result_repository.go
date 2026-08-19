package repository

import (
	"slib.uz/src/core/domain/entity"
)

type NotificationSendResultRepository interface {
	Create(*entity.NotificationSendResultEntity) error
}
