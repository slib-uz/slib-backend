package article_application

import (
	"net/http"

	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/paymentusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type ConfirmPaymentHandler struct {
	uc *ConfirmPaymentUseCase
}

// @inject
func NewConfirmPaymentHandler(uc *ConfirmPaymentUseCase) *ConfirmPaymentHandler {
	return &ConfirmPaymentHandler{uc: uc}
}

// Handle
// @Tags article-application
// @Accept json
// @Produce json
// @Param applicationId path int true "Application ID"
// @Success 200 {object} interface{}
// @Router /article-application/{applicationId}/confirm-payment [post]
func (this *ConfirmPaymentHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	applicationId, err := context.GetIntPathParam(ctx, "applicationId")
	if err != nil {
		return err
	}

	err = this.uc.Execute(uint(applicationId))
	if err != nil {
		return err
	}
	return c.JsonResponse(http.StatusOK, map[string]interface{}{
		"message": "Payment confirmed successfully",
	})
}
