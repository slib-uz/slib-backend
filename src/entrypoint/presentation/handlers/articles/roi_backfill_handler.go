package articles

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/roiusecase"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type ROIBackfillHandler struct {
	uc *roiusecase.ROIBackfillUseCase
}

// @inject
func NewROIBackfillHandler(uc *roiusecase.ROIBackfillUseCase) *ROIBackfillHandler {
	return &ROIBackfillHandler{uc: uc}
}

// Handle godoc
// @Tags         article
// @Accept       json
// @Produce      json
// @Security     BasicAuth
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /articles/roi-backfill [post]
func (this *ROIBackfillHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	result, err := this.uc.Execute()
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
