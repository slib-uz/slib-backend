package journal

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/journalusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type JournalListHandler struct {
	uc *journalusecases.JournalListUseCase
}

// @inject
func NewJournalListHandler(uc *journalusecases.JournalListUseCase) *JournalListHandler {
	return &JournalListHandler{uc: uc}
}

// Handle JournalListHandler
// @Tags         journal
// @Param        page          query uint   false "Page number"
// @Param        page_size     query uint   false "Page size"
// @Param        publisher     query int    false "Publisher ID"
// @Param        name          query string false "Journal name"
// @Param        description   query string false "Journal description"
// @Param from_year query int false "Start year for filtering"
// @Param to_year query int false "End year for filtering"
// @Param        issn          query string false "ISSN"
// @Param        publisher_name query string false "Publisher name"
// @Param        study_field query string false "StudyFieldIDS"
// @Param        language   query string false "LanguageIds"
// @Param        indexing_type query string false "IndexingTypes (WOS, SCOPUS, RINC, SPRINGER, E-LIBRARY, RINTS)"
// @Param        submission_access  query     int     false  "Filter by submission access"  Enums(10,-10)
// @Param        oak          query     bool    false  "Filter by OAK registration (true=registered, false=not registered)"
// @Param        sort_by      query     string  false  "Sort by field"  Enums(views_count,rating_sum,established_date)
// @Param        order        query     string  false  "Sort order"  Enums(asc,desc)
// @Accept       json
// @Produce      json
// @Success      200 {array} entity.JournalEntity
// @Failure      400 {object} response.Response
// @Security     BearerAuth
// @Router       /journal/list [get]
func (this *JournalListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, size := context2.GetPagingParams(c)
	publisherId := this.getPublisherId(c)
	name := this.getSearchParams(c)
	issn := this.getISSNParam(c)
	studyFieldIds, err := context2.GetUintListQueryParam(ctx, "study_field")
	languageIds, err := context2.GetUintListQueryParam(ctx, "language")
	description := this.getDescriptionParam(c)
	publisherName := this.getPublisherNameParam(c)
	fromYear := this.getFromYearParam(c)
	toYear := this.getToYearParam(c)
	submissionAccess := context2.GetIntQueryParam(c, "submission_access", 0)
	oakRegistered := context2.GetBoolQueryParam(c, "oak")
	indexingTypes := context2.GetStringListQueryParam(c, "indexing_type")
	sortBy := c.QueryParam("sort_by")
	order := c.QueryParam("order")

	paging, err := this.uc.Execute(page, size, submissionAccess, oakRegistered, publisherId, name, description, issn, publisherName, languageIds, studyFieldIds, fromYear, toYear, indexingTypes, sortBy, order)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, paging)
}

func (this *JournalListHandler) getPublisherId(ctx echo.Context) uint {
	publisherId, err := strconv.Atoi(ctx.QueryParam("publisher"))
	if err != nil {
		return 0
	}
	id := uint(publisherId)
	return id
}

func (this *JournalListHandler) getDescriptionParam(ctx echo.Context) (description *string) {
	descriptionParam := ctx.QueryParam("description")
	if descriptionParam != "" {
		description = &descriptionParam
	} else {
		description = nil
	}
	return description
}

func (this *JournalListHandler) getFromYearParam(ctx echo.Context) (fromYear *int) {
	fromYearParam := ctx.QueryParam("from_year")
	if fromYearParam == "" {
		return nil
	}

	year, err := strconv.Atoi(fromYearParam)
	if err != nil || year <= 0 {
		return nil
	}
	return &year
}

func (this *JournalListHandler) getToYearParam(ctx echo.Context) (toYear *int) {
	toYearParam := ctx.QueryParam("to_year")
	if toYearParam == "" {
		return nil
	}

	year, err := strconv.Atoi(toYearParam)
	if err != nil || year <= 0 {
		return nil
	}
	return &year
}

func (this *JournalListHandler) getPublisherNameParam(ctx echo.Context) (publisherName *string) {
	publisherNameParam := ctx.QueryParam("publisher_name")
	if publisherNameParam != "" {
		publisherName = &publisherNameParam
	} else {
		publisherName = nil
	}
	return publisherName
}
func (this *JournalListHandler) getSearchParams(ctx echo.Context) (name *string) {
	qname := ctx.QueryParam("name")

	if qname != "" {
		name = &qname
	}
	return name
}

func (this *JournalListHandler) getISSNParam(ctx echo.Context) (issn *string) {
	issnParam := ctx.QueryParam("issn")
	if issnParam != "" {
		issn = &issnParam
	} else {
		issn = nil
	}
	return issn
}
