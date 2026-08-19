package test

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/service"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/notification/schema"
)

type SendToTokensHandler struct {
	service *service.SendTestNotificationUseCase
}

// @inject
func NewSendToTokensHandler(service *service.SendTestNotificationUseCase) *SendToTokensHandler {
	return &SendToTokensHandler{service: service}
}

// Handle godoc
// @Tags notification-test
// @Accept json
// @Produce json
// @Security BasicAuth
// @Param request body schema.SendTestNotificationRequest true "Send Test Notification Request"
// @Success 200 {object} response.Response "Notification sent successfully"
// @Router /notification-test/send-to-tokens [post]
func (this *SendToTokensHandler) Handle(ctx echo.Context) error {

	c := ctx.(*context2.Context)

	data, err := context2.GetBody[schema.SendTestNotificationRequest](ctx)
	if err != nil {
		return err
	}

	result, err := this.service.SendToTokens(data.Tokens, data.Title, data.Body)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
