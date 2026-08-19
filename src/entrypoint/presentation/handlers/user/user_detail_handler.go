package user

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/userusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type UserDetailHandler struct {
	uc *userusecases.UserDetailUseCase
}

// @inject
func NewUserDetailHandler(uc *userusecases.UserDetailUseCase) *UserDetailHandler {
	return &UserDetailHandler{uc: uc}
}

// Handle
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "User ID"
// @Success      200  {object} entity.UserDetailEntity
// @Failure      403  {object} response.Response
// @Router       /user/{id} [get]
func (this *UserDetailHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)
	user := c.User

	if user == nil {
		return echo.ErrForbidden
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.ErrBadRequest
	}

	result, err := this.uc.Execute(user, uint(id))
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
