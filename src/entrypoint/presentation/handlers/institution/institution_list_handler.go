package institution

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/institutionusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type InstitutionListHandler struct {
	uc *InstitutionListUseCase
}

// @inject
func NewInstitutionListHandler(uc *InstitutionListUseCase) *InstitutionListHandler {
	return &InstitutionListHandler{uc: uc}
}

// Handle
// @Tags Institution
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param tin query string false "Search by TIN"
// @Param name query string false "Search by Name"
// @Success 200 {array} entity.InstitutionEntity
// @Router /institution/list [get]
func (this *InstitutionListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, size := context2.GetPagingParams(c)

	tin := c.QueryParam("tin")
	name := c.QueryParam("name")

	result, err := this.uc.Execute(page, size, &tin, &name)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
