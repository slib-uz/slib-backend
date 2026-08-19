package support_dialog

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/support_dialog_usecase"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/support_dialog/schema"
)

type SupportDialogCreateAnswerHandler struct {
	uc *SupportDialogCreateAnswerUseCase
}

// @inject
func NewSupportDialogCreateAnswerHandler(uc *SupportDialogCreateAnswerUseCase) *SupportDialogCreateAnswerHandler {
	return &SupportDialogCreateAnswerHandler{uc: uc}
}

// Handle SupportDialogCreateAnswerHandler
// @Summary Create support answer
// @Description Create support answer
// @Security BearerAuth
// @Tags Support Dialog
// @Accept json
// @Produce json
// @Param request body schema.SupportDialogCreateAnswerRequest true "SupportDialogCreateAnswerRequest"
// @Success 201 {string} string "Your answer has been sent to user"
// @Router /support-dialog-answer/create [post]
func (this SupportDialogCreateAnswerHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	data, err := context2.GetBody[schema.SupportDialogCreateAnswerRequest](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User.ID, data.ToEntity()); err != nil {
		return err
	}

	return c.JsonResponse(201, "Your answer has been sent to user")
}
