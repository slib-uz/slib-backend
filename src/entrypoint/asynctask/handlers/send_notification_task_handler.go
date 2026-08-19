package handlers

import (
	"context"

	"slib.uz/src/core/application/usecase/notificationusecases"
	"slib.uz/src/core/domain/entity"
)

type SendNotificationTaskHandler struct {
	uc *notificationusecases.SendNotificationUseCase
}

// @inject
func NewSendNotificationTaskHandler(uc *notificationusecases.SendNotificationUseCase) *SendNotificationTaskHandler {
	return &SendNotificationTaskHandler{uc: uc}
}

func (h *SendNotificationTaskHandler) ProcessTask(_ context.Context, task *entity.AsyncTask) error {
	n, err := UnmarshalPayload[entity.NotificationEntity](task)
	if err != nil {
		return err
	}
	return h.uc.Execute(&n)
}
