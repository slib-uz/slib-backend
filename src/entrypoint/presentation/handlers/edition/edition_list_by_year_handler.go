package edition

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/editionusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type EditionByYearListHandler struct {
	uc *editionusecases.JournalEditionsListByYearUseCase
}

// @inject
func NewEditionByYearListHandler(uc *editionusecases.JournalEditionsListByYearUseCase) *EditionByYearListHandler {
	return &EditionByYearListHandler{uc: uc}
}

// Handle godoc
// @Tags         edition
// @Accept       json
// @Produce      json
// @Param        journalId  path      int     true   "Journal ID"
// @Param        page       query     int     false  "Page number"
// @Param        page_size  query     int     false  "Items per page"
// @Success      200  {object}  editionusecases.EditionListResponse
// @Failure      400  {object}  response.Response
// @Router       /edition/{journalId}/list-by-year [get]
func (this *EditionByYearListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	journalIdValue, err := context.GetIntPathParam(ctx, "journalId")
	if err != nil {
		return err
	}
	journalID := uint(journalIdValue)

	page, pageSize := context.GetPagingParams(c)

	// Call usecase to get paginated editions grouped by year
	response, err := this.uc.Execute(ctx.Request().Context(), journalID, page, pageSize)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, response)
}
