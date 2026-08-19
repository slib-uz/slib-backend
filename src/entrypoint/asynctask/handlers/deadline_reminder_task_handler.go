package handlers

import (
	"context"

	"slib.uz/src/core/application/usecase/deadlineusecases"
	"slib.uz/src/core/domain/entity"
)

type DeadlineReminderTaskHandler struct {
	uc *deadlineusecases.DeadlineReminderUseCase
}

// @inject
func NewDeadlineReminderTaskHandler(uc *deadlineusecases.DeadlineReminderUseCase) *DeadlineReminderTaskHandler {
	return &DeadlineReminderTaskHandler{uc: uc}
}

func (h *DeadlineReminderTaskHandler) ProcessTask(_ context.Context, _ *entity.AsyncTask) error {
	return h.uc.Execute()
}
