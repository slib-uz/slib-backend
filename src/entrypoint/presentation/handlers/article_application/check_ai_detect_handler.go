package article_application

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/aidetectusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type CheckAiDetectHandler struct {
	uc *CheckAiDetectUseCase
}

// @inject
func NewCheckAiDetectHandler(uc *CheckAiDetectUseCase) *CheckAiDetectHandler {
	return &CheckAiDetectHandler{uc: uc}
}

// Handle
// @Tags         article-application
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        stageId path int true "Stage ID"
// @Success      200 {object} response.Response
// @Router       /article-application/{stageId}/ai-detect/check [post]
func (this *CheckAiDetectHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	stageID, err := context.GetIntPathParam(ctx, "stageId")
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User, uint(stageID)); err != nil {
		return err
	}

	return c.JsonResponse(200, "Submit for AI detection check successfully initiated")
}
