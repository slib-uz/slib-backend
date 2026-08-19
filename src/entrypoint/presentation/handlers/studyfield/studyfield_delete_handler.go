package studyfield

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/studyfieldusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type StudyFieldDeleteHandler struct {
	uc *StudyFieldDeleteUseCase
}

// @inject
func NewStudyFieldDeleteHandler(uc *StudyFieldDeleteUseCase) *StudyFieldDeleteHandler {
	return &StudyFieldDeleteHandler{uc: uc}
}

// Handle StudyFieldDeleteHandler
// @Tags studyfield
// @Accept json
// @Produce json
// @Param id path int true "Study Field ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response "Study field not found"
// @Router /studyfield/manage/delete/{id} [delete]
func (this *StudyFieldDeleteHandler) Handle(ctx echo.Context) error {

	c := ctx.(*context2.Context)

	id, err := context2.GetIntPathParam(ctx, "id")
	if err != nil {
		return err
	}

	if err := this.uc.Execute(uint(id)); err != nil {
		return err
	}
	return c.JsonResponse(200, "Study field deleted successfully")
}
