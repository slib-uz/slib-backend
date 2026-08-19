package handlers

import (
	"context"

	"slib.uz/src/core/application/usecase/logusecases"
	"slib.uz/src/core/domain/entity"
)

type SendTelegramAlertTaskHandler struct {
	uc *logusecases.SendTelegramAlertUseCase
}

// @inject
func NewSendTelegramAlertTaskHandler(uc *logusecases.SendTelegramAlertUseCase) *SendTelegramAlertTaskHandler {
	return &SendTelegramAlertTaskHandler{uc: uc}
}

func (h *SendTelegramAlertTaskHandler) ProcessTask(_ context.Context, task *entity.AsyncTask) error {
	alert, err := UnmarshalPayload[entity.TelegramAlertEntity](task)
	if err != nil {
		return err
	}
	return h.uc.Execute(alert.Message)
}
