package about

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/aboutusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type AboutRetrieveHandler struct {
	uc *aboutusecases.AboutDetailUseCase
}

// @inject
func NewAboutRetrieveHandler(uc *aboutusecases.AboutDetailUseCase) *AboutRetrieveHandler {
	return &AboutRetrieveHandler{uc: uc}
}

// Handle
// @Tags         About
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "About ID"
// @Success      200  {object}  entity.AboutEntity
// @Router       /about/detail/{id} [get]
func (this *AboutRetrieveHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest
	}

	res, err := this.uc.Execute(uint(id))
	if err != nil {
		return err
	}
	if res == nil {
		return echo.ErrNotFound
	}
	return c.JsonResponse(200, res)
}
