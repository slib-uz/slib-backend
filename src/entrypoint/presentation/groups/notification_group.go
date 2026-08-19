package groups

import (
	"github.com/labstack/echo/v4"
	notification2 "slib.uz/src/entrypoint/presentation/handlers/notification"
)

//import (
//	container "slib.uz/src/di"
//	"slib.uz/src/presentation/app"
//	"slib.uz/src/presentation/interceptor/permissions"
//)
//
//var NotificationGroup = app.GetInstance().Group(
//	"/notification",
//	container.InitJwtAuthMiddleware().Call,
//	permissions.AuthenticatedPermission,
//)

//func init() {
//	NotificationGroup.GET("/unread-count", container.InitUnreadCountHandler().Handle)
//	NotificationGroup.POST("/register-token", container.InitRegisterTokenHandler().Handle)
//	NotificationGroup.DELETE("/remove-token/:token", container.InitRemoveTokenHandler().Handle)
//	NotificationGroup.GET("/my/list", container.InitMyNotificationsListHandler().Handle)
//	NotificationGroup.GET("/my/detail/:id", container.InitNotificationDetailHandler().Handle)
//}

type NotificationGroup struct {
	unreadCountHandler         *notification2.UnreadCountHandler
	registerTokenHandler       *notification2.RegisterTokenHandler
	removeTokenHandler         *notification2.RemoveTokenHandler
	myNotificationsListHandler *notification2.MyNotificationsListHandler
	notificationDetailHandler  *notification2.NotificationDetailHandler
}

// @inject
func NewNotificationGroup(unreadCountHandler *notification2.UnreadCountHandler, registerTokenHandler *notification2.RegisterTokenHandler, removeTokenHandler *notification2.RemoveTokenHandler, myNotificationsListHandler *notification2.MyNotificationsListHandler, notificationDetailHandler *notification2.NotificationDetailHandler) *NotificationGroup {
	return &NotificationGroup{unreadCountHandler: unreadCountHandler, registerTokenHandler: registerTokenHandler, removeTokenHandler: removeTokenHandler, myNotificationsListHandler: myNotificationsListHandler, notificationDetailHandler: notificationDetailHandler}
}

func (this *NotificationGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/unread-count", this.unreadCountHandler.Handle)
	group.POST("/register-token", this.registerTokenHandler.Handle)
	group.DELETE("/remove-token/:token", this.removeTokenHandler.Handle)
	group.GET("/my/list", this.myNotificationsListHandler.Handle)
	group.GET("/my/detail/:id", this.notificationDetailHandler.Handle)
}
