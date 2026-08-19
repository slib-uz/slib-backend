package studyfield

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/studyfieldusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/studyfield/schema"
)

type StudyFieldUpdateHandler struct {
	uc *StudyFieldUpdateUseCase
}

// @inject
func NewStudyFieldUpdateHandler(uc *StudyFieldUpdateUseCase) *StudyFieldUpdateHandler {
	return &StudyFieldUpdateHandler{uc: uc}
}

// Handle StudyFieldUpdateHandler
// @Tags studyfield
// @Accept json
// @Produce json
// @Param studyField body schema.StudyFieldUpdateRequest true "Study Field"
// @Success 200 {object} response.Response
// @Router /studyfield/manage/update [patch]
func (this *StudyFieldUpdateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	studyField, err := context2.GetBody[schema.StudyFieldUpdateRequest](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(studyField.ToEntity()); err != nil {
		return err
	}

	return c.JsonResponse(200, "Study field updated successfully")
}
