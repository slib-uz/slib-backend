package journal

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/journalusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/journal/schema"
)

type RemoveReviewersHandler struct {
	uc *RemoveReviewersUsecase
}

// @inject
func NewRemoveReviewersHandler(uc *RemoveReviewersUsecase) *RemoveReviewersHandler {
	return &RemoveReviewersHandler{uc: uc}
}

// Handle
// @Tags journal-manage
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body schema.RemoveReviewerRequest true "Remove reviewers from a journal"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /journal-manage/reviewers/remove [post]
func (this *RemoveReviewersHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	data, err := context2.GetBody[schema.RemoveReviewerRequest](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(data.JournalID, data.ReviewerID); err != nil {
		return err
	}
	return c.JsonResponse(200, nil)
}
