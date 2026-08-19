package doi

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/doiusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type CrossRefDepositListHandler struct {
	uc *doiusecases.CrossRefDepositListUseCase
}

// @inject
func NewCrossRefDepositListHandler(uc *doiusecases.CrossRefDepositListUseCase) *CrossRefDepositListHandler {
	return &CrossRefDepositListHandler{uc: uc}
}

// Handle
// @Tags         doi-settings
// @Security     BearerAuth
// @Produce      json
// @Param        page      query int false "Page number"
// @Param        page_size query int false "Page size"
// @Param        journal_id query int false "Journal ID"
// @Success      200 {array} entity.DoiDepositEntity
// @Failure      400,401,403,500 {object} response.Response
// @Router       /doi-settings/deposits [get]
func (this *CrossRefDepositListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)
	page, pageSize := context.GetPagingParams(c)
	journalID := uint(context.GetIntQueryParam(c, "journal_id", 0))

	result, err := this.uc.Execute(c.Request().Context(), c.User, journalID, page, pageSize)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
