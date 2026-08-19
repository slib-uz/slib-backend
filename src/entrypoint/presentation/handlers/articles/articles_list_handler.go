package articles

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/articleusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type ArticlesListHandler struct {
	uc *ArticlesListUseCase
}

// @inject
func NewArticlesListHandler(uc *ArticlesListUseCase) *ArticlesListHandler {
	return &ArticlesListHandler{uc: uc}
}

// Handle
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        journal      query     int     false  "Filter by journal ID"
// @Param        language     query     string  false  "Filter by language IDs (comma-separated)"
// @Param        studyfields  query     string  false  "Filter by study field"
// @Param        tags         query     string  false "Filter by tags"
// @Param        article_title query     string  false  "Search by article title"
// @Param        author_name   query     string  false  "Search by author name"
// @Param        journal_name  query     string  false  "Search by journal name"
// @Param        issn_paper    query     string  false  "Search by ISSN paper"
// @Param        issn_online   query     string  false  "Search by ISSN online"
// @Param        access_type  query     int     false  "Filter by access type"  Enums(10,-10)
// @Param        sort         query     string  false  "Sort by field"  Enums(views_count,rating_sum,publication_date)
// @Param        page         query     int false  "Page number for pagination"
// @Param        page_size    query     int false  "Number of items per page for pagination"
// @Success      200  {array} entity.ArticleInputEntity
// @Failure      400  {object}  response.Response
// @Param        has_doi      query     bool    false  "Filter by DOI presence (true=has DOI, false=no DOI)"
// @Param        from_year    query     int     false  "Start year for filtering"
// @Param        to_year      query     int     false  "End year for filtering"
// @Router       /articles/list [get]
func (this *ArticlesListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	fromYear := this.GetYearParam(c, "from_year")
	toYear := this.GetYearParam(c, "to_year")
	journal := context.GetIntQueryParam(c, "journal", 0)
	languages, err := context.GetUintListQueryParam(c, "language")
	if err != nil {
		return echo.NewHTTPError(400, "Invalid language parameter")
	}
	studyfields, err := context.GetUintListQueryParam(c, "studyfields")
	if err != nil {
		return echo.NewHTTPError(400, "Invalid studyfields parameter")
	}
	tags := context.GetStringListQueryParam(c, "tags")

	articleTitle := this.getStringParam(c, "article_title")
	authorName := this.getStringParam(c, "author_name")
	journalName := this.getStringParam(c, "journal_name")
	issnPaper := this.getStringParam(c, "issn_paper")
	issnOnline := this.getStringParam(c, "issn_online")
	sort := c.QueryParam("sort")

	page, pageSize := context.GetPagingParams(c)
	accessType := context.GetIntQueryParam(c, "access_type", 0)
	hasDoi := context.GetBoolQueryParam(c, "has_doi")
	result, err := this.uc.Execute(fromYear, toYear, uint(journal), languages, studyfields, tags, articleTitle, authorName, journalName, issnPaper, issnOnline, sort, page, pageSize, accessType, hasDoi)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}

func (this *ArticlesListHandler) GetYearParam(c echo.Context, param string) *int {
	year := context.GetIntQueryParam(c, param, 0)
	if year <= 0 {
		return nil
	}
	return &year
}

func (this *ArticlesListHandler) getStringParam(ctx echo.Context, paramName string) *string {
	value := ctx.QueryParam(paramName)
	if value == "" {
		return nil
	}
	return &value
}
