package articles

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	. "slib.uz/src/core/application/usecase/articleusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type UserArticlesListHandler struct {
	uc *UserArticleListUseCase
}

// @inject
func NewUserArticlesListHandler(uc *UserArticleListUseCase) *UserArticlesListHandler {
	return &UserArticlesListHandler{uc: uc}
}

// Handle
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        scienceId  path      string  true  "Author Science ID"
// @Param        page               query     int     false "Page number" default(1)
// @Param        page_size          query     int     false "Page size" default(10)
// @Param        search             query     string  false "Search query"
// @Param 		 study_field 		query 	  string  false	"StudyField"
// @Param 		 year				query	  string  false "Year"
// @Param        journal_id         query     int     false "Filter by journal ID (from current journal domain/tenant)"
// @Success      200  {array} entity.ArticleInputEntity
// @Failure      400  {object}  response.Response
// @Router       /articles/author/{scienceId} [get]
func (this *UserArticlesListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(ctx)
	searchParam := ctx.QueryParam("search")
	studyFieldIds, err := context2.GetUintListQueryParam(ctx, "study_field")
	year := context2.GetIntQueryParam(ctx, "year", 0)

	var journalID *uint
	if val := c.QueryParam("journal_id"); val != "" {
		id, parseErr := strconv.ParseUint(val, 10, 64)
		if parseErr != nil {
			return response.NewFailResponse(400, "Invalid journal_id")
		}
		rid := uint(id)
		journalID = &rid
	}

	paging, err := this.uc.Execute(ctx.Param("scienceId"), page, pageSize, searchParam, studyFieldIds, year, journalID)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, paging)
}
