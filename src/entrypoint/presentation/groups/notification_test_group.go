package groups

import (
	"github.com/labstack/echo/v4"
	test2 "slib.uz/src/entrypoint/presentation/handlers/notification/test"
)

type NotificationTestGroup struct {
	sendToTokensHandler *test2.SendToTokensHandler
	sendToTopicHandler  *test2.SendToTopicHandler
	errorAlertHandler   *test2.ErrorAlertHandler
}

// @inject
func NewNotificationTestGroup(sendToTokensHandler *test2.SendToTokensHandler, sendToTopicHandler *test2.SendToTopicHandler, errorAlertHandler *test2.ErrorAlertHandler) *NotificationTestGroup {
	return &NotificationTestGroup{sendToTokensHandler: sendToTokensHandler, sendToTopicHandler: sendToTopicHandler, errorAlertHandler: errorAlertHandler}
}

func (this *NotificationTestGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/send-to-tokens", this.sendToTokensHandler.Handle)
	group.POST("/send-to-topic", this.sendToTopicHandler.Handle)
	group.GET("/error-alert", this.errorAlertHandler.Handle)
}
