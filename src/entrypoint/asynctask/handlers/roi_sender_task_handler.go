package handlers

import (
	"context"
	"log"

	"slib.uz/src/core/application/tasks"
	"slib.uz/src/core/application/usecase/roiusecase"
	"slib.uz/src/core/domain/entity"
)

type RoiSenderTaskHandler struct {
	uc *roiusecase.ArticlePublishROIUseCase
}

// @inject
func NewRoiSenderTaskHandler(uc *roiusecase.ArticlePublishROIUseCase) *RoiSenderTaskHandler {
	return &RoiSenderTaskHandler{uc: uc}
}

func (h *RoiSenderTaskHandler) ProcessTask(_ context.Context, task *entity.AsyncTask) error {
	payload, err := UnmarshalPayload[tasks.SetRoiSenderPayload](task)
	if err != nil {
		log.Println("[RoiSenderTaskHandler] UnmarshalPayload:", err.Error())
		return err
	}

	if _, err := h.uc.Execute(payload.User, payload.ApplicationID, payload.File); err != nil {
		log.Println("[RoiSenderTaskHandler] Invoke:", err.Error())
		return err
	}

	return nil
}
