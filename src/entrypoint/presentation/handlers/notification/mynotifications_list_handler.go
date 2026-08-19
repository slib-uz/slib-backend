package notification

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/notificationusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type MyNotificationsListHandler struct {
	uc *UserNotificationsListUseCase
}

// @inject
func NewMyNotificationsListHandler(uc *UserNotificationsListUseCase) *MyNotificationsListHandler {
	return &MyNotificationsListHandler{uc: uc}
}

// Handle godoc
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param unread query bool false "Filter unread notifications" default(false)
// @Success 200 {array} entity.NotificationEntity "List of user notifications"
// @Router /notification/my/list [get]
func (this *MyNotificationsListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)
	page, pageSize := context2.GetPagingParams(c)

	paging, err := this.uc.Execute(c.User.ID, this.unreadParam(ctx), page, pageSize)
	if err != nil {
		return err
	}
	return c.JsonResponse(200, paging)
}

func (this *MyNotificationsListHandler) unreadParam(ctx echo.Context) bool {
	unread := ctx.QueryParam("unread")
	return unread == "true"
}
