package journalstatistics

import (
	"github.com/labstack/echo/v4"
	journalstatisticsusecase "slib.uz/src/core/application/usecase/journal_statistics_usecase"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type JournalStatisticsListExcelHandler struct {
	uc *journalstatisticsusecase.JournalStatisticListExcelUseCase
}

// @inject
func NewJournalStatisticsListExcelHandler(uc *journalstatisticsusecase.JournalStatisticListExcelUseCase) *JournalStatisticsListExcelHandler {
	return &JournalStatisticsListExcelHandler{uc: uc}
}

// Handle
// @Tags journal-statistics
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param institution_id query int false "Institution ID"
// @Success 200 {file} file
// @Failure 500 {object} response.Response
// @Router /journal-statistics/list/excel [get]
func (this *JournalStatisticsListExcelHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	institutionID := uint(context2.GetIntQueryParam(c, "institution_id", 0))

	file, err := this.uc.Execute(institutionID)
	if err != nil {
		return err
	}

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return err
	}

	c.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=journal_statistics.xlsx")

	return c.Blob(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
}
