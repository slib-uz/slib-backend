package article_application

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/article_applications_usecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type ApplicationsListHandler struct {
	uc *ApplicationsListUseCase
}

// @inject
func NewApplicationsListHandler(uc *ApplicationsListUseCase) *ApplicationsListHandler {
	return &ApplicationsListHandler{uc: uc}
}

// Handle ApplicationsListHandler
// @Tags article-application
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param start_date query string false "Start date" format(date) default(2023-01-01)
// @Param end_date query string false "End date" format(date) default(2023-12-31)
// @Param search query string false "Search query"
// @Param journalId query int false "Journal ID"
// @Param status query int false "Status"
// @Success 200 {array} entity.ApplicationBasicEntity
// @Router /article-application/list [get]
func (this *ApplicationsListHandler) Handle(ctx echo.Context) error {

	c := ctx.(*context2.Context)

	page, size := context2.GetPagingParams(ctx)
	startDate, endDate := context2.GetDateRangeQueryParams(ctx, "start_date", "end_date")
	q := c.QueryParam("search")
	journalId := uint(context2.GetIntQueryParam(c, "journalId", 0))
	status := context2.GetIntQueryParam(c, "status", -999)

	paging, err := this.uc.Execute(c.User, journalId, page, size, startDate, endDate, q, status)
	if err != nil {
		return err
	}
	return c.JsonResponse(200, paging)
}
