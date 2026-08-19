package myapplication

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/myapplicationusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type MyApplicationsListHandler struct {
	uc *MyApplicationsListUseCase
}

// @inject
func NewMyApplicationsListHandler(uc *MyApplicationsListUseCase) *MyApplicationsListHandler {
	return &MyApplicationsListHandler{uc: uc}
}

// Handle
// @Tags my-application
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param journal_id query int false "Journal ID filter"
// @Success 200 {array} entity.ApplicationBasicEntity
// @Router /my-application/list [get]
func (this *MyApplicationsListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(c)
	journalID := context2.GetIntQueryParam(c, "journal_id", 0)

	paging, err := this.uc.Execute(c.User.ID, page, pageSize, uint(journalID))
	if err != nil {
		return err
	}

	return c.JsonResponse(200, paging)
}
