package handlers

import (
	"context"

	"slib.uz/src/core/application/usecase/journalusecases"
	"slib.uz/src/core/domain/entity"
)

type JournalViewsCountUpdateTaskHandler struct {
	uc *journalusecases.UpdateJournalViewsCountUseCase
}

// @inject
func NewJournalViewsCountUpdateTaskHandler(uc *journalusecases.UpdateJournalViewsCountUseCase) *JournalViewsCountUpdateTaskHandler {
	return &JournalViewsCountUpdateTaskHandler{uc: uc}
}

func (h *JournalViewsCountUpdateTaskHandler) ProcessTask(_ context.Context, _ *entity.AsyncTask) error {
	return h.uc.Execute()
}
