package partner

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/partnerusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type PartnersListHandler struct {
	uc *PartnerListUseCase
}

// @inject
func NewPartnersListHandler(uc *PartnerListUseCase) *PartnersListHandler {
	return &PartnersListHandler{uc: uc}
}

// Handle
// @Tags partner
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {array} entity.PartnerEntity
// @Router /partner/list [get]
func (this *PartnersListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(ctx)
	paging, err := this.uc.Execute(page, pageSize)

	if err != nil {
		return err
	}

	return c.JsonResponse(200, paging)
}
