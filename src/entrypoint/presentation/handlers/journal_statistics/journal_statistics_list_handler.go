package journalstatistics

import (
	"github.com/labstack/echo/v4"
	journalstatisticsusecase "slib.uz/src/core/application/usecase/journal_statistics_usecase"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type JournalStatisticsListHandler struct {
	uc *journalstatisticsusecase.JournalStatisticListUseCase
}

// @inject
func NewJournalStatisticsListHandler(uc *journalstatisticsusecase.JournalStatisticListUseCase) *JournalStatisticsListHandler {
	return &JournalStatisticsListHandler{uc: uc}
}

// Handle
// @Tags journal-statistics
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param institution_id query int false "Institution ID"
// @Param publisher_id query int false "Publisher ID"
// @Param name query string false "Journal name"
// @Param description query string false "Journal description"
// @Param issn query string false "ISSN"
// @Param publisher_name query string false "Publisher name"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /journal-statistics/list [get]
func (this *JournalStatisticsListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(ctx)
	institutionID := uint(context2.GetIntQueryParam(c, "institution_id", 0))
	publisherID := uint(context2.GetIntQueryParam(c, "publisher_id", 0))
	name := context2.GetStringQueryParam(c, "name")
	description := context2.GetStringQueryParam(c, "description")
	issn := context2.GetStringQueryParam(c, "issn")
	publisherName := context2.GetStringQueryParam(c, "publisher_name")

	result, err := this.uc.Execute(page, pageSize, institutionID, publisherID, name, description, issn, publisherName)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
