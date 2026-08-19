package entity

import (
	"encoding/json"
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type OutboxEventEntity struct {
	ID          uint             `json:"id"`
	EventID     string           `json:"event_id"`
	Version     int64            `json:"version"`
	Event       enum.OutboxEvent `json:"event"`
	Payload     json.RawMessage  `json:"payload"`
	DeliveredAt *time.Time       `json:"delivered_at"`
}

func NewOutboxEventEntity(eventID string, version int64, event enum.OutboxEvent, payload json.RawMessage) *OutboxEventEntity {
	return &OutboxEventEntity{
		EventID: eventID,
		Version: version,
		Event:   event,
		Payload: payload,
	}
}
