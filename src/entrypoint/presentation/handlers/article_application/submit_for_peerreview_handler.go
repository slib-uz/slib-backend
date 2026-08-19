package article_application

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/etaqrizusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/article_application/schema"
	"slib.uz/src/entrypoint/presentation/headers"
)

type SubmitForPeerReviewHandler struct {
	uc *SubmitForPeerReviewUseCase
}

// @inject
func NewSubmitForPeerReviewHandler(uc *SubmitForPeerReviewUseCase) *SubmitForPeerReviewHandler {
	return &SubmitForPeerReviewHandler{uc: uc}
}

// Handle
// @Tags article-application
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body schema.SubmitForPeerReviewRequest true "Submit for peer review request"
// @Param X-Idempotency-Key header string true "Idempotency Key"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /article-application/peer-review/submit [post]
func (this *SubmitForPeerReviewHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	idempotencyKey := ctx.Request().Header.Get(headers.HeaderIdempotencyKey)

	data, err := context2.GetBody[schema.SubmitForPeerReviewRequest](ctx)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(idempotencyKey, c.User.ID, data.ApplicationID, data.ReviewerIDs); err != nil {
		return err
	}

	return c.JsonResponse(200, map[string]any{"message": "Successfully submitted for peer review"})
}
