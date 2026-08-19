package article_application

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/antiplagusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type CheckAntiPlagHandler struct {
	uc *CheckAntiPlagUseCase
}

// @inject
func NewCheckAntiPlagHandler(uc *CheckAntiPlagUseCase) *CheckAntiPlagHandler {
	return &CheckAntiPlagHandler{uc: uc}
}

// Handle
// @Tags article-application
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param stageId path int true "Application ID"
// @Success 200 {object} response.Response
// @Router /article-application/{stageId}/antiplag/check [post]
func (this *CheckAntiPlagHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	stageId, err := context2.GetIntPathParam(ctx, "stageId")
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User, uint(stageId)); err != nil {
		return err
	}

	return c.JsonResponse(200, "Submit for anti-plagiarism check successfully initiated")
}
