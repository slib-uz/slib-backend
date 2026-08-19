package studyfield

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/studyfieldusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/studyfield/schema"
)

type StudyFieldCreateHandler struct {
	uc *StudyFieldCreateUseCase
}

// @inject
func NewStudyFieldCreateHandler(uc *StudyFieldCreateUseCase) *StudyFieldCreateHandler {
	return &StudyFieldCreateHandler{uc: uc}
}

// Handle
// @Tags studyfield
// @Accept json
// @Produce json
// @Param studyField body schema.StudyFieldCreateRequest true "Study Field"
// @Success 200 {object} response.Response
// @Router /studyfield/manage/create [post]
func (this *StudyFieldCreateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	body, err := context2.GetBody[schema.StudyFieldCreateRequest](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(body.ToEntity()); err != nil {
		return err
	}

	return c.JsonResponse(201, "Study field created successfully")
}
