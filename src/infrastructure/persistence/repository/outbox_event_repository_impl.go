package repository

import (
	"time"

	"gorm.io/datatypes"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/persistence/models"
)

type OutboxEventRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewOutboxEventRepository(baseRepository *BaseRepository) repository.OutboxEventRepository {
	return &OutboxEventRepositoryImpl{BaseRepository: baseRepository}
}

func (this *OutboxEventRepositoryImpl) Create(event *entity.OutboxEventEntity) (*entity.OutboxEventEntity, error) {
	model := &models.OutboxEventModel{
		EventID: event.EventID,
		Version: event.Version,
		Event:   event.Event,
		Payload: datatypes.JSON(event.Payload),
	}

	if err := this.db().Create(model).Error; err != nil {
		return nil, err
	}

	event.ID = model.ID
	return event, nil
}

func (this *OutboxEventRepositoryImpl) MarkDelivered(eventID string) error {
	now := time.Now()
	return this.db().Model(&models.OutboxEventModel{}).
		Where("event_id = ?", eventID).
		Update("delivered_at", &now).Error
}
