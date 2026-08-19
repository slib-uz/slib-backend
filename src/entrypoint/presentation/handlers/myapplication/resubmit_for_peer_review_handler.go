package myapplication

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/etaqrizusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/myapplication/schema"
)

type ResubmitForPeerReviewHandler struct {
	uc *ReSubmitForPeerReviewUseCase
}

// @inject
func NewResubmitForPeerReviewHandler(uc *ReSubmitForPeerReviewUseCase) *ResubmitForPeerReviewHandler {
	return &ResubmitForPeerReviewHandler{uc: uc}
}

// Handle
// @Security BearerAuth
// @Tags my-application
// @Accept json
// @Produce json
// @Param request body schema.ResubmitForPeerReviewRequest true "Resubmit for peer review request"
// @Success 200 {object} response.Response
// @Router /my-application/resubmit-for-peer-review [post]
func (this *ResubmitForPeerReviewHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	data, err := context.GetBody[schema.ResubmitForPeerReviewRequest](ctx)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User.ID, data.SubmissionID, data.File); err != nil {
		return err
	}

	return c.JsonResponse(200, "Application resubmitted successfully")
}
