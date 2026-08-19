package site

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/statisticsusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type SiteStatisticsHandler struct {
	uc *statisticsusecases.SiteStatisticsUseCase
}

// @inject
func NewSiteStatisticsHandler(uc *statisticsusecases.SiteStatisticsUseCase) *SiteStatisticsHandler {
	return &SiteStatisticsHandler{uc: uc}
}

// Handle godoc
// @Tags Site
// @Accept json
// @Produce json
// @Param journal_id query int false "Filter by journal ID"
// @Param publisher_id query int false "Filter by publisher ID"
// @Param institution_id query int false "Filter by institution ID"
// @Success 200 {object} entity.SiteStatisticsEntity
// @Router /site/statistic [get]
func (this *SiteStatisticsHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	var journalID *uint
	if val := ctx.QueryParam("journal_id"); val != "" {
		id, _ := strconv.ParseUint(val, 10, 64)
		rid := uint(id)
		journalID = &rid
	}

	var publisherID *uint
	if val := ctx.QueryParam("publisher_id"); val != "" {
		id, _ := strconv.ParseUint(val, 10, 64)
		pid := uint(id)
		publisherID = &pid
	}

	var institutionID *uint
	if val := ctx.QueryParam("institution_id"); val != "" {
		id, _ := strconv.ParseUint(val, 10, 64)
		iid := uint(id)
		institutionID = &iid
	}

	stats, err := this.uc.Execute(journalID, publisherID, institutionID)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, stats)
}
