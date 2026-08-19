package userfind

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/userusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type UserFindHandler struct {
	uc *UserFindUseCase
}

// @inject
func NewUserFindHandler(uc *UserFindUseCase) *UserFindHandler {
	return &UserFindHandler{uc: uc}
}

// Handle UserFindHandler
// @Tags user-find
// @Accept  json
// @Produce  json
// @Param scienceId query string true "Science ID"
// @Success 200 {object} entity.UserBasicEntity
// @Security BearerAuth
// @Failure 404 {object} response.Response
// @Router /users/find [get]
func (this *UserFindHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)
	user, err := this.uc.Execute(c.User, ctx.QueryParam("scienceId"))
	if err != nil {
		return err
	}
	return c.JsonResponse(200, user)
}
