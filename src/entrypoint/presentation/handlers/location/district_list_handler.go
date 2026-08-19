package location

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/locationusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type DistrictListHandler struct {
	uc *locationusecases.DistrictListUseCase
}

// @inject
func NewDistrictListHandler(uc *locationusecases.DistrictListUseCase) *DistrictListHandler {
	return &DistrictListHandler{uc: uc}
}

// Handle godoc
// @Tags         location
// @Accept       json
// @Produce      json
// @Param        region_id  query  int  true  "Region ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /location/districts/list [get]
func (this *DistrictListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	regionID := context2.GetIntQueryParam(c, "region_id", 0)
	if regionID <= 0 {
		return response.NewFailResponse(400, "region_id is required")
	}

	districts, err := this.uc.Execute(uint(regionID))
	if err != nil {
		return err
	}

	return c.JsonResponse(200, districts)
}
