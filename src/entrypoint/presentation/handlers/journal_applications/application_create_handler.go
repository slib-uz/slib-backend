package journal_applications

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/journal_applications_usecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type JournalApplicationCreateHandler struct {
	uc *journal_applications_usecases.JournalCreateUseCase
}

// @inject
func NewJournalApplicationCreateHandler(uc *journal_applications_usecases.JournalCreateUseCase) *JournalApplicationCreateHandler {
	return &JournalApplicationCreateHandler{uc: uc}
}

// @Tags journal-applications
// @Accept json
// @Produce json
// @Param request body entity.JournalCreateEntity true "Journal"
// @Success 201 {object} map[string]any
// @Security BearerAuth
// @Router /journal-applications/create [post]
func (this *JournalApplicationCreateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	request, err := context.GetBody[entity.JournalCreateEntity](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(request); err != nil {
		return err
	}

	return c.JsonResponse(http.StatusCreated, "Journal created successfully")
}
