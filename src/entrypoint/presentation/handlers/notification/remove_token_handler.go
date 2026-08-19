package notification

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/notificationusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type RemoveTokenHandler struct {
	uc *RemoveTokenUseCase
}

// @inject
func NewRemoveTokenHandler(uc *RemoveTokenUseCase) *RemoveTokenHandler {
	return &RemoveTokenHandler{uc: uc}
}

// Handle godoc
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param token path string true "Notification Token"
// @Success 200 {object} response.Response "Token removed successfully"
// @Router /notification/remove-token/{token} [delete]
func (this *RemoveTokenHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	token := c.Param("token")

	if err := this.uc.Execute(c.User.ID, token); err != nil {
		return err
	}
	return c.JsonResponse(200, "Token removed successfully")
}
