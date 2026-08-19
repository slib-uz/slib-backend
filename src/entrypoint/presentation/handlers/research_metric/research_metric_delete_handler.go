package research_metric

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/researchmetricusecase"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type ResearchMetricDeleteHandler struct {
	uc *ResearchMetricDeleteUseCase
}

// @inject
func NewResearchMetricDeleteHandler(uc *ResearchMetricDeleteUseCase) *ResearchMetricDeleteHandler {
	return &ResearchMetricDeleteHandler{uc: uc}
}

// Handle ResearchMetricDeleteHandler
// @Tags research-metric
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Research Metric ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /research-metric/{id} [delete]
func (this *ResearchMetricDeleteHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)
	id, err := context2.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}
	if err := this.uc.Execute(c.User.ID, uint(id)); err != nil {
		return c.JsonResponse(400, err.Error())
	}
	return c.JsonResponse(200, "Research metric deleted successfully")
}
