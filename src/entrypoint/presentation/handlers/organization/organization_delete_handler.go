package organization

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/organizationusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type OrganizationDeleteHandler struct {
	uc *organizationusecases.OrganizationDeleteUseCase
}

// @inject
func NewOrganizationDeleteHandler(uc *organizationusecases.OrganizationDeleteUseCase) *OrganizationDeleteHandler {
	return &OrganizationDeleteHandler{uc: uc}
}

// Handle
// @Tags Organization
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Organization ID"
// @Success 200 {string} string "Success"
// @Failure 403 {string} string "Forbidden"
// @Router /organization/delete/{id} [delete]
func (this *OrganizationDeleteHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	id, err := context2.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User, uint(id)); err != nil {
		return err
	}

	return c.JsonResponse(200, "Success")
}
