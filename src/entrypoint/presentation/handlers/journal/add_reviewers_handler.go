package journal

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/journalusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/journal/schema"
)

type AddReviewersHandler struct {
	uc *AddReviewersUsecase
}

// @inject
func NewAddReviewersHandler(uc *AddReviewersUsecase) *AddReviewersHandler {
	return &AddReviewersHandler{uc: uc}
}

// Handle
// @Tags journal-manage
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body schema.AddReviewersRequest true "Add Reviewers Request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /journal-manage/reviewers/add [post]
func (this *AddReviewersHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	data, err := context2.GetBody[schema.AddReviewersRequest](ctx)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(data.JournalID, data.ReviewerIds); err != nil {
		return err
	}

	return c.JsonResponse(200, nil)
}
