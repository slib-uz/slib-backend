package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity/enum"
)

type NotificationTemplateModel struct {
	gorm.Model
	Key   enum.NotificationTemplate `gorm:"size:64;not null;uniqueIndex"`
	Title datatypes.JSON
	Body  datatypes.JSON
}
