package repository

import "slib.uz/src/core/domain/entity"

type OutboxEventRepository interface {
	Create(event *entity.OutboxEventEntity) (*entity.OutboxEventEntity, error)
	MarkDelivered(eventID string) error
}
