package groups

import (
	"github.com/labstack/echo/v4"
	research_metric2 "slib.uz/src/entrypoint/presentation/handlers/research_metric"
)

type ResearchMetricGroup struct {
	createOrUpdateHandler *research_metric2.ResearchMetricCreateOrUpdateHandler
	deleteHandler         *research_metric2.ResearchMetricDeleteHandler
}

// @inject
func NewResearchMetricGroup(createOrUpdateHandler *research_metric2.ResearchMetricCreateOrUpdateHandler, deleteHandler *research_metric2.ResearchMetricDeleteHandler) *ResearchMetricGroup {
	return &ResearchMetricGroup{
		createOrUpdateHandler: createOrUpdateHandler,
		deleteHandler:         deleteHandler,
	}
}

func (this *ResearchMetricGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/create-or-update", this.createOrUpdateHandler.Handle)
	group.DELETE("/:id", this.deleteHandler.Handle)
}
