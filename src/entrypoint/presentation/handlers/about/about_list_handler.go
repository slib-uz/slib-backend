package about

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/aboutusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type AboutListHandler struct {
	uc *aboutusecases.AboutListUseCase
}

// @inject
func NewAboutListHandler(uc *aboutusecases.AboutListUseCase) *AboutListHandler {
	return &AboutListHandler{uc: uc}
}

// Handle
// @Tags         About
// @Accept       json
// @Produce      json
// @Success      200  {array}   entity.AboutEntity
// @Router       /about/list [get]
func (this *AboutListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	res, err := this.uc.Execute()
	if err != nil {
		return err
	}
	return c.JsonResponse(200, res)
}
