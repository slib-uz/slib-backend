package role

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/roleusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/role/schema"
)

type UserRoleCreateHandler struct {
	uc *roleusecases.UserRoleCreateUseCase
}

// @inject
func NewUserRoleCreateHandler(uc *roleusecases.UserRoleCreateUseCase) *UserRoleCreateHandler {
	return &UserRoleCreateHandler{uc: uc}
}

// Handle UserRoleCreateHandler
// @Tags role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Role body schema.RoleCreateRequest true "Role"
// @Success 200 {object} response.Response
// @Router /role/create [post]
func (this *UserRoleCreateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	data, err := context2.GetBody[schema.RoleCreateRequest](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User, data.ToEntity()); err != nil {
		return err
	}
	return c.JsonResponse(http.StatusCreated, "Role created successfully")
}
