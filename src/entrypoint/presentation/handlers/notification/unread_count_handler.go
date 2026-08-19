package notification

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/notificationusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type UnreadCountHandler struct {
	uc *UnreadCountUseCase
}

// @inject
func NewUnreadCountHandler(uc *UnreadCountUseCase) *UnreadCountHandler {
	return &UnreadCountHandler{uc: uc}
}

// Handle godoc
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "Unread notification count"
// @Router /notification/unread-count [get]
func (this *UnreadCountHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	count, err := this.uc.Execute(c.User.ID)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, map[string]any{"count": count})
}
