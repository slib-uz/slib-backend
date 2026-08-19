package outboxusecases

import (
	"github.com/labstack/gommon/log"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
)

type DeliveryOutboxEventUseCase struct {
	billingGateway gateway.SlibBillingGateway
	outboxRepo     repository.OutboxEventRepository
}

// @inject
func NewDeliveryOutboxEventUseCase(
	billingGateway gateway.SlibBillingGateway,
	outboxRepo repository.OutboxEventRepository,
) *DeliveryOutboxEventUseCase {
	return &DeliveryOutboxEventUseCase{
		billingGateway: billingGateway,
		outboxRepo:     outboxRepo,
	}
}

func (uc *DeliveryOutboxEventUseCase) Execute(event *entity.OutboxEventEntity) error {
	if err := uc.billingGateway.SendEvent(event); err != nil {
		log.Printf("Error sending outbox event %s: %v\n", event.EventID, err)
		return err
	}

	if err := uc.outboxRepo.MarkDelivered(event.EventID); err != nil {
		log.Printf("Error marking outbox event as delivered: %v\n", err)
	}

	return nil
}
