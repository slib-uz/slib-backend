package notification

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/notificationusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/notification/schema"
)

type RegisterTokenHandler struct {
	uc *RegisterTokenUseCase
}

// @inject
func NewRegisterTokenHandler(uc *RegisterTokenUseCase) *RegisterTokenHandler {
	return &RegisterTokenHandler{uc: uc}
}

// Handle godoc
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body schema.RegisterTokenRequest true "Register Token Request"
// @Success 200 {object} response.Response "Token registered successfully"
// @Router /notification/register-token [post]
func (this *RegisterTokenHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)
	data, err := context2.GetBody[schema.RegisterTokenRequest](ctx)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User.ID, data.Token, data.Segment); err != nil {
		return err
	}

	return c.JsonResponse(200, "Token registered successfully")
}
