package ethics_policy

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/ethics_policy_usecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type EthicsPolicyDetailHandler struct {
	uc *EthicsPolicyDetailUseCase
}

// @inject
func NewEthicsPolicyDetailHandler(uc *EthicsPolicyDetailUseCase) *EthicsPolicyDetailHandler {
	return &EthicsPolicyDetailHandler{uc: uc}
}

// Handle
// @Summary      Get list of ethics policies
// @Tags         ethics_policy
// @Accept       json
// @Produce      json
// @Param        id      path     int     true     "Ethics policy id"
// @Success      200  {object}   entity.EthicsPolicyEntity
// @Failure      400  {object}  response.Response
// @Router       /ethics_policy/detail/{id} [get]
func (this *EthicsPolicyDetailHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	id, err := context.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}

	ethics_policy, err := this.uc.Execute(uint(id))

	if err != nil {
		return err
	}

	return c.JsonResponse(200, ethics_policy)
}
