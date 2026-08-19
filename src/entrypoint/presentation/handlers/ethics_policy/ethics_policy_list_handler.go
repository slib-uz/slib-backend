package ethics_policy

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/ethics_policy_usecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type EthicsPolicyListHandler struct {
	uc *EthicsPolicyListUseCase
}

// @inject
func NewEthicsPolicyListHandler(uc *EthicsPolicyListUseCase) *EthicsPolicyListHandler {
	return &EthicsPolicyListHandler{uc: uc}
}

// Handle
// @Summary      Get list of ethics policies
// @Tags         ethics_policy
// @Accept       json
// @Produce      json
// @Param        page      query     int     false     "Page number"
// @Param        page_size query     int     false     "Page size"
// @Success      200  {array}   entity.EthicsPolicyEntity
// @Failure      400  {object}  response.Response
// @Router       /ethics_policy/list [get]
func (this *EthicsPolicyListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	page, pageSize := context.GetPagingParams(ctx)
	paging, err := this.uc.Execute(page, pageSize)

	if err != nil {
		return err
	}

	return c.JsonResponse(200, paging)
}
