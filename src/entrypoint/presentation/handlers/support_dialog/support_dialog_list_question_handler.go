package support_dialog

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/support_dialog_usecase"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type SupportDialogQuestionListHandler struct {
	uc *SupportDialogQuestionListUseCase
}

// @inject
func NewSupportDialogQuestionListHandler(uc *SupportDialogQuestionListUseCase) *SupportDialogQuestionListHandler {
	return &SupportDialogQuestionListHandler{uc: uc}
}

// Handle SupportDialogQuestionListHandler
// @Summary List support dialogs
// @Description List support dialogs
// @Security BearerAuth
// @Tags Support Dialog
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param page_size query int false "Page Size"
// @Param ordering query string false "Ordering" Enums(created_at,-created_at)
// @Success 200 {array} entity.SupportDialogEntity
// @Failure 400 {object} response.Response
// @Router /support-dialog-question/list [get]
func (this SupportDialogQuestionListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(ctx)
	ordering := ctx.QueryParam("ordering")
	userID := c.User.ID
	paging, err := this.uc.Execute(page, pageSize, ordering, userID)

	if err != nil {
		return err
	}

	return c.JsonResponse(200, paging)
}
