package doi

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/doiusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/doi/schema"
)

type DoiSettingUpdateHandler struct {
	uc *doiusecases.DoiSettingUpdateUseCase
}

// @inject
func NewDoiSettingUpdateHandler(uc *doiusecases.DoiSettingUpdateUseCase) *DoiSettingUpdateHandler {
	return &DoiSettingUpdateHandler{uc: uc}
}

// Handle
// @Tags         doi-settings
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        journalId path int true "Journal ID"
// @Param        request body schema.DoiSettingRequest true "DOI Setting"
// @Success      200 {object} map[string]any
// @Failure      400,401,403,404 {object} response.Response
// @Router       /doi-settings/{journalId} [put]
func (this *DoiSettingUpdateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	journalID, err := context.GetIntPathParam(c, "journalId")
	if err != nil {
		return err
	}

	body, err := context.GetBody[schema.DoiSettingRequest](c)
	if err != nil {
		return err
	}

	e := body.ToEntity(uint(journalID))
	if err := this.uc.Execute(c.Request().Context(), e); err != nil {
		return err
	}

	return c.JsonResponse(200, "DOI setting updated successfully")
}
