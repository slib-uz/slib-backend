package article_application

import (
	"net/http"

	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/article_applications_usecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type ApplicationCheckForPaymentHandler struct {
	uc *ApplicationCheckForPaymentUsecase
}

// @inject
func NewApplicationCheckForPaymentHandler(uc *ApplicationCheckForPaymentUsecase) *ApplicationCheckForPaymentHandler {
	return &ApplicationCheckForPaymentHandler{uc: uc}
}

// Handle
// @Tags article-application
// @Accept json
// @Produce json
// @Param applicationId path int true "Application ID"
// @Success 200 {object} entity.ReviewStageEntity
// @Security BearerAuth
// @Router /article-application/{applicationId}/check-for-payment [get]
func (this *ApplicationCheckForPaymentHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	applicationId, err := context.GetIntPathParam(ctx, "applicationId")
	if err != nil {
		return err
	}

	stage, err := this.uc.Execute(uint(applicationId))
	if err != nil {
		return err
	}
	return c.JsonResponse(http.StatusOK, stage)
}
