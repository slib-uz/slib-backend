package guide

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/guideusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type GuidesListHandler struct {
	uc *GuideListUseCase
}

// @inject
func NewGuideListHandler(uc *GuideListUseCase) *GuidesListHandler {
	return &GuidesListHandler{uc: uc}
}

// Handle
// @Summary      Get list of guides
// @Tags         guide
// @Accept       json
// @Produce      json
// @Param        page      query     int     false     "Page number"
// @Param        page_size query     int     false     "Page size"
// @Success      200  {array}   entity.GuideListEntity
// @Failure      400  {object}  response.Response
// @Router       /guides/list [get]
func (this *GuidesListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(ctx)
	paging, err := this.uc.Execute(page, pageSize)

	if err != nil {
		return err
	}

	return c.JsonResponse(200, paging)
}
