package handlers

import (
	"context"

	"slib.uz/src/core/application/usecase/outboxusecases"
	"slib.uz/src/core/domain/entity"
)

type OutboxEventTaskHandler struct {
	uc *outboxusecases.DeliveryOutboxEventUseCase
}

// @inject
func NewOutboxEventTaskHandler(uc *outboxusecases.DeliveryOutboxEventUseCase) *OutboxEventTaskHandler {
	return &OutboxEventTaskHandler{uc: uc}
}

func (h *OutboxEventTaskHandler) ProcessTask(_ context.Context, task *entity.AsyncTask) error {
	event, err := UnmarshalPayload[entity.OutboxEventEntity](task)
	if err != nil {
		return err
	}
	return h.uc.Execute(&event)
}
