package location

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/locationusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type RegionListHandler struct {
	uc *locationusecases.RegionListUseCase
}

// @inject
func NewRegionListHandler(uc *locationusecases.RegionListUseCase) *RegionListHandler {
	return &RegionListHandler{uc: uc}
}

// Handle godoc
// @Tags         location
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response
// @Router       /location/regions/list [get]
func (this *RegionListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	regions, err := this.uc.Execute()
	if err != nil {
		return err
	}

	return c.JsonResponse(200, regions)
}
