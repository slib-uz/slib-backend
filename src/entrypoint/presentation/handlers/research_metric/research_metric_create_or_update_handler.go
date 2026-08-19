package research_metric

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/researchmetricusecase"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/research_metric/schema"
)

type ResearchMetricCreateOrUpdateHandler struct {
	uc *ResearchMetricCreateOrUpdateUseCase
}

// @inject
func NewResearchMetricCreateOrUpdateHandler(uc *ResearchMetricCreateOrUpdateUseCase) *ResearchMetricCreateOrUpdateHandler {
	return &ResearchMetricCreateOrUpdateHandler{uc: uc}
}

// Handle ResearchMetricCreateOrUpdateHandler
// @Tags research-metric
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param research-metric body schema.ResearchMetricCreateOrUpdateRequest true "ResearchMetricCreateOrUpdateRequest"
// @Success 200 {object} response.Response
// @Router /research-metric/create-or-update [post]
func (this *ResearchMetricCreateOrUpdateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	data, err := context2.GetBody[schema.ResearchMetricCreateOrUpdateRequest](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User.ID, data.ToEntity()); err != nil {
		return c.JsonResponse(400, err.Error())
	}
	return c.JsonResponse(200, "Research metric created or updated successfully")
}
