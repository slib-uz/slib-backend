package role

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/roleusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type RoleDeleteHandler struct {
	uc *RoleDeleteUseCase
}

// @inject
func NewRoleDeleteHandler(uc *RoleDeleteUseCase) *RoleDeleteHandler {
	return &RoleDeleteHandler{uc: uc}
}

// Handle RoleDeleteHandler
// @Tags role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /role/delete/{id} [delete]
func (this *RoleDeleteHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)
	id, err := context2.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}

	err = this.uc.Execute(uint(id))
	if err != nil {
		return err
	}
	return c.JsonResponse(200, "Role deleted successfully")
}
