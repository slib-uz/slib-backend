package support_dialog

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/support_dialog_usecase"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type SupportDialogAnswerListHandler struct {
	uc *SupportDialogAnswerListUseCase
}

// @inject
func NewSupportDialogAnswerListHandler(uc *SupportDialogAnswerListUseCase) *SupportDialogAnswerListHandler {
	return &SupportDialogAnswerListHandler{uc: uc}
}

// Handle SupportDialogAnswerListHandler
// @Summary List support answers
// @Description List support answers
// @Security BearerAuth
// @Tags Support Dialog
// @Accept json
// @Produce json
// @Param chatID path int true "Chat ID"
// @Param page query int false "Page"
// @Param page_size query int false "Page Size"
// @Param ordering query string false "Ordering" Enums(created_at,-created_at)
// @Success 200 {array} entity.SupportDialogEntity
// @Failure 400 {object} response.Response
// @Router /support-dialog-answer/list/{chatID} [get]
func (this SupportDialogAnswerListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(ctx)
	ordering := ctx.QueryParam("ordering")

	chatId, err := context2.GetIntPathParam(c, "chatID")
	if err != nil {
		return err
	}

	paging, err := this.uc.Execute(page, pageSize, ordering, uint(chatId))

	if err != nil {
		return err
	}

	return c.JsonResponse(200, paging)
}
