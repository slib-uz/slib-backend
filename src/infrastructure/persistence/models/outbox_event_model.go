package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity/enum"
)

type OutboxEventModel struct {
	gorm.Model
	EventID     string           `gorm:"not null;uniqueIndex;size:50"`
	Version     int64            `gorm:"not null"`
	Event       enum.OutboxEvent `gorm:"not null;size:50"`
	Payload     datatypes.JSON
	DeliveredAt *time.Time
}

func (*OutboxEventModel) TableName() string {
	return "outbox_events"
}
