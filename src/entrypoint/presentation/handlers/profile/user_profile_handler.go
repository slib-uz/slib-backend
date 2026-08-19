package profile

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/userusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type ProfileHandler struct {
	usecases *userusecases.UserProfileUseCase
}

// @inject
func NewProfileHandler(usecases *userusecases.UserProfileUseCase) *ProfileHandler {
	return &ProfileHandler{usecases: usecases}
}

// Handle ProfileHandler
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.UserProfileEntity
// @Security BearerAuth
// @Router /user/profile [get]
func (this *ProfileHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	data, err := this.usecases.Execute(c.User, c.User.ID)

	if err != nil {
		return err
	}

	return c.JsonResponse(200, data)
}
