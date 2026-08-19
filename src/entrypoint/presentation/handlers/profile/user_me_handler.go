package profile

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/userusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type UserMeHandler struct {
	uc *userusecases.UserMeUseCase
}

// @inject
func NewUserMeHandler(uc *userusecases.UserMeUseCase) *UserMeHandler {
	return &UserMeHandler{uc: uc}
}

// Handle UserMeHandler
// @Tags User
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} entity.UserMeEntity
// @Router /user/me [get]
func (this *UserMeHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	user, err := this.uc.Execute(c.User.ID)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, user)
}
