package article_application

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/aidetectusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type AiDetectResultListHandler struct {
	uc *AiDetectResultListUseCase
}

// @inject
func NewAiDetectResultListHandler(uc *AiDetectResultListUseCase) *AiDetectResultListHandler {
	return &AiDetectResultListHandler{uc: uc}
}

// Handle
// @Tags article-application
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param applicationId path int true "Application ID"
// @Param reviewStageId query int true "Review Stage ID"
// @Success 200 {array} entity.AiDetectResultEntity
// @Router /article-application/{applicationId}/ai-detect/results [get]
func (this *AiDetectResultListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	applicationID, err := context2.GetIntPathParam(ctx, "applicationId")
	if err != nil {
		return err
	}

	results, err := this.uc.Execute(
		c.User,
		uint(applicationID),
		uint(context2.GetIntQueryParam(ctx, "reviewStageId", 0)),
	)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, results)
}
