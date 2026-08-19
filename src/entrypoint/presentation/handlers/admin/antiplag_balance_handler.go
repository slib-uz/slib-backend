package admin

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/antiplagusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type AntiPlagBalanceHandler struct {
	uc *BalanceUseCase
}

// @inject
func NewAntiPlagBalanceHandler(uc *BalanceUseCase) *AntiPlagBalanceHandler {
	return &AntiPlagBalanceHandler{uc: uc}
}

// Handle godoc
// @Summary      Get AntiPlag Balance
// @Description  Get AntiPlag Balance
// @Tags         admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any
// @Router 	 /admin/antiplag/balance [get]
func (this *AntiPlagBalanceHandler) Handle(ctx echo.Context) error {

	c := ctx.(*context.Context)

	balance, err := this.uc.Execute()
	if err != nil {
		return err
	}

	return c.JsonResponse(200, map[string]any{"balance": balance})
}
