package notification

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/notificationusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type NotificationDetailHandler struct {
	uc *NotificationDetailUseCase
}

// @inject
func NewNotificationDetailHandler(uc *NotificationDetailUseCase) *NotificationDetailHandler {
	return &NotificationDetailHandler{uc: uc}
}

// Handle godoc
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Notification ID"
// @Success 200 {object} entity.NotificationEntity "Notification details"
// @Router /notification/my/detail/{id} [get]
func (this *NotificationDetailHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	id, err := context2.GetIntPathParam(ctx, "id")
	if err != nil {
		return err
	}

	n, err := this.uc.Execute(c.User.ID, uint(id))
	if err != nil {
		return err
	}

	return c.JsonResponse(200, n)
}
