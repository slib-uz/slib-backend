package articles

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/articleusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type JournalArticlesListHandler struct {
	uc *JournalArticlesListUseCase
}

// @inject
func NewJournalArticlesListHandler(uc *JournalArticlesListUseCase) *JournalArticlesListHandler {
	return &JournalArticlesListHandler{uc: uc}
}

// Handle godoc
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        journalId  path      int  true  "Journal ID"
// @Param        page               query     int     false "Page number" default(1)
// @Param        page_size          query     int     false "Page size" default(10)
// @Param        search             query     string  false "Search query"
// @Param 		 author				query 	  string  false "Author"
// @Param 		 study_field 		query 	  string  false	"StudyField"
// @Param 		 year				query	  string  false "Year"
// @Param 		 tag				query	  string  false "Tag"
// @Success      200  {array} entity.ArticleInputEntity
// @Router 	 /articles/journal/{journalId} [get]
func (this *JournalArticlesListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(ctx)

	journalID, err := context2.GetIntPathParam(ctx, "journalId")
	if err != nil {
		return err
	}

	searchParam := ctx.QueryParam("search")
	authorSearch := ctx.QueryParam("author")
	studyFieldIds, err := context2.GetUintListQueryParam(ctx, "study_field")
	year := context2.GetIntQueryParam(ctx, "year", 0)
	tag := ctx.QueryParam("tag")

	paging, err := this.uc.Execute(uint(journalID), page, pageSize, searchParam, authorSearch, tag, studyFieldIds, year)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, paging)

}
