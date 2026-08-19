package journal_applications

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/journal_applications_usecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type JournalApplicationListHandler struct {
	uc *ApplicationsListUsecase
}

// @inject
func NewApplicationListHandler(uc *ApplicationsListUsecase) *JournalApplicationListHandler {
	return &JournalApplicationListHandler{uc: uc}
}

// Handle
// @Tags journal-applications
// @Accept: application/json
// @Produce: application/json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param publisher_id query int false "Publisher ID"
// @Param status query int false "Application status"
// @Success 200 {array} entity.JournalApplicationEntity
// @Failure 400 {object} response.Response
// @Router /journal-applications/applications [get]
func (this *JournalApplicationListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(ctx)
	publisherID := uint(context2.GetIntQueryParam(c, "publisher_id", 0))
	status := context2.GetIntQueryParam(ctx, "status", -100)

	paging, err := this.uc.Execute(c.User, publisherID, page, pageSize, status)

	if err != nil {
		return err
	}

	return c.JsonResponse(200, paging)
}
