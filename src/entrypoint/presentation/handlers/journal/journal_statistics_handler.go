package journal

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/journalusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type JournalStatisticsListHandler struct {
	uc *JournalStatisticsListUseCase
}

// @inject
func NewJournalStatisticsListHandler(uc *JournalStatisticsListUseCase) *JournalStatisticsListHandler {
	return &JournalStatisticsListHandler{uc: uc}
}

// Handle
// @Tags         journal
// @Accept       json
// @Produce      json
// @Param        page         query     int false  "Page number for pagination"
// @Param        page_size    query     int false  "Number of items per page for pagination"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /journal/statistics [get]
func (this *JournalStatisticsListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(c)

	result, err := this.uc.Execute(page, pageSize)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}

type JournalStatisticsExcelHandler struct {
	uc *JournalStatisticsExcelUseCase
}

// @inject
func NewJournalStatisticsExcelHandler(uc *JournalStatisticsExcelUseCase) *JournalStatisticsExcelHandler {
	return &JournalStatisticsExcelHandler{uc: uc}
}

// Handle
// @Tags         journal
// @Accept       json
// @Produce      application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Success      200  {file}  file
// @Failure      500  {object}  response.Response
// @Router       /journal/statistics/excel [get]
func (this *JournalStatisticsExcelHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	data, err := this.uc.Execute()
	if err != nil {
		return err
	}

	c.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=journal_statistics.xlsx")

	return c.Blob(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}
