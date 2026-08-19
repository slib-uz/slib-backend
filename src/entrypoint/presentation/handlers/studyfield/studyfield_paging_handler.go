package studyfield

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/studyfieldusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type StudyFieldPagingHandler struct {
	uc *StudyFieldPagingUseCase
}

// @inject
func NewStudyFieldPagingHandler(uc *StudyFieldPagingUseCase) *StudyFieldPagingHandler {
	return &StudyFieldPagingHandler{uc: uc}
}

// Handle godoc
// @Summary Get study fields with paging
// @Description Get study fields with paging
// @Tags studyfield
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Success 200 {array} entity.StudyFieldEntity
// @Router /studyfield/paging [get]
func (this *StudyFieldPagingHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(ctx)
	search := c.QueryParam("search")

	paging, err := this.uc.Execute(page, pageSize, search)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, paging)

}
