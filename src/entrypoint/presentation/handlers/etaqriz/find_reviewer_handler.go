package etaqriz

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/etaqrizusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type FindReviewerHandler struct {
	uc *FindReviewerUseCase
}

// @inject
func NewFindReviewerHandler(uc *FindReviewerUseCase) *FindReviewerHandler {
	return &FindReviewerHandler{uc: uc}
}

// Handle
// @Tags e-taqriz
// @Accept json
// @Produce json
// @Success 200 {object} entity.ReviewerEntity
// @Failure 404 {object} response.Response
// @Param scienceId query string true "Science ID"
// @Security BearerAuth
// @Router /etaqriz/reviewer/find [get]
func (this *FindReviewerHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	reviewer, err := this.uc.Execute(c.QueryParam("scienceId"))
	if err != nil {
		return err
	}

	return c.JsonResponse(200, reviewer)
}
