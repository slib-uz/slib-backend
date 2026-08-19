package organization

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/organizationusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type OrganizationListHandler struct {
	uc *OrganizationListUseCase
}

// @inject
func NewOrganizationListHandler(uc *OrganizationListUseCase) *OrganizationListHandler {
	return &OrganizationListHandler{uc: uc}
}

// Handle
// @Tags Organization
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param tin query string false "Search by TIN"
// @Param name query string false "Search by Name"
// @Param address query string false "Search by Address"
// @Success 200 {array} entity.OrganizationEntity
// @Router /organization/list [get]
func (this *OrganizationListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, size := context2.GetPagingParams(c)

	tin := c.QueryParam("tin")
	name := c.QueryParam("name")
	address := c.QueryParam("address")

	result, err := this.uc.Execute(page, size, &tin, &name, &address)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
