package handlers

import (
	"context"

	"slib.uz/src/core/application/usecase/newsusecases"
	"slib.uz/src/core/domain/entity"
)

type NewsViewsCountUpdateTaskHandler struct {
	uc *newsusecases.UpdateNewsViewsCountUseCase
}

// @inject
func NewNewsViewsCountUpdateTaskHandler(uc *newsusecases.UpdateNewsViewsCountUseCase) *NewsViewsCountUpdateTaskHandler {
	return &NewsViewsCountUpdateTaskHandler{uc: uc}
}

func (h *NewsViewsCountUpdateTaskHandler) ProcessTask(_ context.Context, _ *entity.AsyncTask) error {
	return h.uc.Execute()
}
