package report

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/reportuscases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type ReportListHandler struct {
	uc *reportuscases.ReportListUseCase
}

// @inject
func NewReportListHandler(uc *reportuscases.ReportListUseCase) *ReportListHandler {
	return &ReportListHandler{uc: uc}
}

// Handle ReportListHandler
// @Tags report
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param ordering query string false "Ordering" Enums(created_at,-created_at)
// @Param target_type query string false "Target type" Enums(journal, article)
// @Success 200 {array} entity.ReportEntity
// @Failure 400 {object} response.Response
// @Router /report/list [get]
func (this ReportListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(ctx)
	ordering := ctx.QueryParam("ordering")
	targetType := ctx.QueryParam("target_type")
	paging, err := this.uc.Execute(page, pageSize, ordering, targetType)
	if err != nil {
		return err
	}
	return c.JsonResponse(200, paging)
}
