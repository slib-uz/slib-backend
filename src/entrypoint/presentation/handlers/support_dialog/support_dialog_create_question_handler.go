package support_dialog

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/support_dialog_usecase"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/support_dialog/schema"
)

type SupportDialogCreateQuestionHandler struct {
	uc *SupportDialogCreateQuestionUseCase
}

// @inject
func NewSupportDialogCreateQuestionHandler(uc *SupportDialogCreateQuestionUseCase) *SupportDialogCreateQuestionHandler {
	return &SupportDialogCreateQuestionHandler{uc: uc}
}

// Handle SupportDialogCreateQuestionHandler
// @Summary Create support question
// @Description Create support question
// @Security BearerAuth
// @Tags Support Dialog
// @Accept json
// @Produce json
// @Param request body schema.SupportDialogCreateQuestionRequest true "SupportDialogCreateQuestionRequest"
// @Success 201 {string} string "Your request has been sent to support"
// @Router /support-dialog-question/create [post]
func (this SupportDialogCreateQuestionHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	data, err := context2.GetBody[schema.SupportDialogCreateQuestionRequest](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User.ID, data.ToEntity()); err != nil {
		return err
	}

	return c.JsonResponse(201, "Your request has been sent to support")
}
