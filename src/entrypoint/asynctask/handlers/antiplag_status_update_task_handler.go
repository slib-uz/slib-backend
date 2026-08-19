package handlers

import (
	"context"

	"slib.uz/src/core/application/usecase/antiplagusecases"
	"slib.uz/src/core/domain/entity"
)

type AntiPlagStatusUpdateTaskHandler struct {
	uc *antiplagusecases.AntiPlagStatusUpdateUseCase
}

// @inject
func NewAntiPlagStatusUpdateTaskHandler(uc *antiplagusecases.AntiPlagStatusUpdateUseCase) *AntiPlagStatusUpdateTaskHandler {
	return &AntiPlagStatusUpdateTaskHandler{uc: uc}
}

func (h *AntiPlagStatusUpdateTaskHandler) ProcessTask(_ context.Context, task *entity.AsyncTask) error {
	externalID, err := UnmarshalPayload[uint](task)
	if err != nil {
		return err
	}
	return h.uc.Execute(externalID)
}
