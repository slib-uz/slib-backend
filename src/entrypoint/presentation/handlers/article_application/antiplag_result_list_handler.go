package article_application

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/antiplagusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type AntiPlagResultListHandler struct {
	uc *AntiPlagResultListUseCase
}

// @inject
func NewAntiPlagResultListHandler(uc *AntiPlagResultListUseCase) *AntiPlagResultListHandler {
	return &AntiPlagResultListHandler{uc: uc}
}

// Handle
// @Tags article-application
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param applicationId path int true "Application ID"
// @Param reviewStageId query int false "Review Stage ID"
// @Success 200 {array} entity.AntiPlagResultEntity
// @Router /article-application/{applicationId}/antiplag/results [get]
func (this *AntiPlagResultListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	applicationId, err := context2.GetIntPathParam(ctx, "applicationId")
	if err != nil {
		return err
	}

	results, err := this.uc.Execute(
		c.User,
		uint(applicationId),
		uint(context2.GetIntQueryParam(ctx, "reviewStageId", 0)),
	)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, results)
}
