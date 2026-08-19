package handlers

import (
	"context"

	"slib.uz/src/core/application/usecase/aidetectusecases"
	"slib.uz/src/core/domain/entity"
)

type AiDetectStatusUpdateTaskHandler struct {
	uc *aidetectusecases.AiDetectStatusUpdateUseCase
}

// @inject
func NewAiDetectStatusUpdateTaskHandler(uc *aidetectusecases.AiDetectStatusUpdateUseCase) *AiDetectStatusUpdateTaskHandler {
	return &AiDetectStatusUpdateTaskHandler{uc: uc}
}

func (h *AiDetectStatusUpdateTaskHandler) ProcessTask(_ context.Context, task *entity.AsyncTask) error {
	externalID, err := UnmarshalPayload[uint](task)
	if err != nil {
		return err
	}
	return h.uc.Execute(externalID)
}
