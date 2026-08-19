package doi

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/doiusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/doi/schema"
)

type DoiSettingGetHandler struct {
	uc *doiusecases.DoiSettingGetUseCase
}

// @inject
func NewDoiSettingGetHandler(uc *doiusecases.DoiSettingGetUseCase) *DoiSettingGetHandler {
	return &DoiSettingGetHandler{uc: uc}
}

// Handle
// @Tags         doi-settings
// @Security     BearerAuth
// @Produce      json
// @Param        journalId path int true "Journal ID"
// @Success      200 {object} schema.DoiSettingResponse
// @Failure      400,401,403,404 {object} response.Response
// @Router       /doi-settings/{journalId} [get]
func (this *DoiSettingGetHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	journalID, err := context.GetIntPathParam(c, "journalId")
	if err != nil {
		return err
	}

	result, err := this.uc.Execute(c.Request().Context(), uint(journalID))
	if err != nil {
		return err
	}

	return c.JsonResponse(200, schema.NewDoiSettingResponse(result))
}
