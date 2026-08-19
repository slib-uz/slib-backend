package journal_applications

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/journal_applications_usecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type JournalApplicationDetailHandler struct {
	uc *ApplicationDetailUsecase
}

// @inject
func NewJournalApplicationDetailHandler(uc *ApplicationDetailUsecase) *JournalApplicationDetailHandler {
	return &JournalApplicationDetailHandler{uc: uc}
}

// Handle
// @Tags journal-applications
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path uint true "Application ID"
// @Success 200 {object} entity.JournalApplicationEntity
// @Failure 404 {object} response.Response
// @Router /journal-applications/applications/{id} [get]
func (this *JournalApplicationDetailHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	id, err := context2.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}

	app, err := this.uc.Execute(uint(id))
	if err != nil {
		return err
	}
	return c.JsonResponse(200, app)
}
