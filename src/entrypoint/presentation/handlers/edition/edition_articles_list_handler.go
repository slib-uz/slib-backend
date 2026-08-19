package edition

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/articleusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type EditionArticlesListHandler struct {
	uc *articleusecases.EditionArticlesListUseCase
}

// @inject
func NewEditionArticlesListHandler(uc *articleusecases.EditionArticlesListUseCase) *EditionArticlesListHandler {
	return &EditionArticlesListHandler{uc: uc}
}

// Handle godoc
// @Tags         edition
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        journal_id    query     int     true   "Journal ID"
// @Param        edition_id    query     int     false  "Edition ID"
// @Param        unassigned    query     bool    false  "Filter articles with no edition (edition_id IS NULL)"
// @Param        page          query     int     false  "Page number"   default(1)
// @Param        page_size     query     int     false  "Items per page" default(10)
// @Success      200  {array}  entity.ArticleInputEntity
// @Failure      400  {object}  response.Response
// @Router       /edition/articles [get]
func (this *EditionArticlesListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	journalID := uint(context.GetIntQueryParam(ctx, "journal_id", 0))
	editionID := uint(context.GetIntQueryParam(ctx, "edition_id", 0))
	unassigned := ctx.QueryParam("unassigned") == "true"
	page, pageSize := context.GetPagingParams(c)

	result, err := this.uc.Execute(journalID, editionID, unassigned, page, pageSize)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
