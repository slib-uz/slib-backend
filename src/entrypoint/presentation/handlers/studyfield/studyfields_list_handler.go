package studyfield

import (
	"strconv"

	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/studyfieldusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type StudyFieldListHandler struct {
	uc *StudyFieldListUseCase
}

// @inject
func NewStudyFieldListHandler(uc *StudyFieldListUseCase) *StudyFieldListHandler {
	return &StudyFieldListHandler{uc: uc}
}

// Handle StudyFieldListHandler
// @Tags studyfield
// @Accept json
// @Produce json
// @Param journal_id query int false "Filter by journal ID"
// @Param search query string false "Search by name (uz, ru, en)"
// @Success 200 {array} entity.StudyFieldEntity
// @Router /studyfield/list [get]
func (this StudyFieldListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	var journalID *uint
	journalIDParam := c.QueryParam("journal_id")
	if journalIDParam != "" {
		if id, err := strconv.ParseUint(journalIDParam, 10, 32); err == nil {
			idUint := uint(id)
			journalID = &idUint
		}
	}

	search := c.QueryParam("search")

	list, err := this.uc.Execute(journalID, search)
	if err != nil {
		return err
	}
	return c.JsonResponse(200, list)
}
