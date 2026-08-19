package repository

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type NotificationTemplateRepository interface {
	GetByKey(key enum.NotificationTemplate) (*entity.NotificationTemplateEntity, error)
}
