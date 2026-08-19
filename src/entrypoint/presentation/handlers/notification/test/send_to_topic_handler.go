package test

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/service"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/notification/schema"
)

type SendToTopicHandler struct {
	service *service.SendTestNotificationUseCase
}

// @inject
func NewSendToTopicHandler(service *service.SendTestNotificationUseCase) *SendToTopicHandler {
	return &SendToTopicHandler{service: service}
}

// Handle godoc
// @Tags notification-test
// @Accept json
// @Produce json
// @Security BasicAuth
// @Param request body schema.SendTestNotificationRequest true "Send Test Notification Request"
// @Success 200 {object} response.Response "Notification sent successfully"
// @Router /notification-test/send-to-topic [post]
func (this *SendToTopicHandler) Handle(ctx echo.Context) error {
	data, err := context.GetBody[schema.SendTestNotificationRequest](ctx)
	if err != nil {
		return err
	}

	if err := this.service.SendToTopic(data.Topic, data.Title, data.Body); err != nil {
		return err
	}

	return ctx.JSON(200, nil)
}
