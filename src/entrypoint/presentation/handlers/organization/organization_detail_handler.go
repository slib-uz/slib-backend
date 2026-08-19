package organization

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/organizationusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type OrganizationDetailHandler struct {
	uc *OrganizationDetailUseCase
}

// @inject
func NewOrganizationDetailHandler(uc *OrganizationDetailUseCase) *OrganizationDetailHandler {
	return &OrganizationDetailHandler{uc: uc}
}

// Handle
// @Tags Organization
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} entity.OrganizationEntity
// @Failure 404 {object} response.Response
// @Router /organization/detail/{id} [get]
func (this *OrganizationDetailHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	id, err := context2.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}

	organization, err := this.uc.Execute(uint(id))
	if err != nil {
		return err
	}

	return c.JsonResponse(200, organization)
}
