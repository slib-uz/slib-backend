package role

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/roleusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type UserRoleListHandler struct {
	uc *roleusecases.UserRoleAllUseCase
}

// @inject
func NewUserRoleListHandler(uc *roleusecases.UserRoleAllUseCase) *UserRoleListHandler {
	return &UserRoleListHandler{uc: uc}
}

// Handle UserRoleListHandler
// @Tags role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} entity.UserRoleEntity
// @Router /role/user-roles [get]
func (this *UserRoleListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	useRoles, err := this.uc.Execute(c.User.ID)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, useRoles)
}
