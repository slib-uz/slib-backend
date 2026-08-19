package journalstatistics

import (
	"github.com/labstack/echo/v4"
	journalstatisticsusecase "slib.uz/src/core/application/usecase/journal_statistics_usecase"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type JournalStatisticsDetailHandler struct {
	uc *journalstatisticsusecase.JournalStatisticsDetailUseCase
}

// @inject
func NewJournalStatisticsDetailHandler(uc *journalstatisticsusecase.JournalStatisticsDetailUseCase) *JournalStatisticsDetailHandler {
	return &JournalStatisticsDetailHandler{uc: uc}
}

// Handle
// @Tags journal-statistics
// @Accept json
// @Produce json
// @Param journal_id path int true "Journal ID"
// @Success 200 {object} entity.JournalStatisticV2Entity
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /journal-statistics/detail/{journal_id} [get]
func (this *JournalStatisticsDetailHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	journalIDValue, err := context2.GetIntPathParam(c, "journal_id")
	if err != nil {
		return err
	}

	result, err := this.uc.Execute(uint(journalIDValue))
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
