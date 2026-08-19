package edition

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/editionusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type EditionListHandler struct {
	uc *editionusecases.JournalEditionsListUseCase
}

// @inject
func NewEditionListHandler(uc *editionusecases.JournalEditionsListUseCase) *EditionListHandler {
	return &EditionListHandler{uc: uc}
}

// Handle godoc
// @Tags         edition
// @Accept       json
// @Produce      json
// @Param        journalId  path      int     true   "Journal ID"
// @Param        page       query     int     false  "Page number"
// @Param        page_size  query     int     false  "Items per page"
// @Param        search     query     string  false  "Search by name"
// @Param        year       query     int     false  "Filter by year"
// @Success      200  {array}  entity.EditionEntity
// @Failure      400  {object}  response.Response
// @Router       /edition/{journalId}/list [get]
func (this *EditionListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	journalIdValue, err := context.GetIntPathParam(ctx, "journalId")
	if err != nil {
		return err
	}
	journalID := uint(journalIdValue)

	page, pageSize := context.GetPagingParams(c)
	search := ctx.QueryParam("search")
	year := context.GetIntQueryParam(ctx, "year", 0)

	result, err := this.uc.Execute(ctx.Request().Context(), journalID, page, pageSize, search, year)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
