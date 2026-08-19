package doi

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/doiusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/doi/schema"
)

type DoiSettingCreateHandler struct {
	uc *doiusecases.DoiSettingCreateUseCase
}

// @inject
func NewDoiSettingCreateHandler(uc *doiusecases.DoiSettingCreateUseCase) *DoiSettingCreateHandler {
	return &DoiSettingCreateHandler{uc: uc}
}

// Handle
// @Tags         doi-settings
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        journalId path int true "Journal ID"
// @Param        request body schema.DoiSettingRequest true "DOI Setting"
// @Success      201 {object} map[string]any
// @Failure      400,401,403,409 {object} response.Response
// @Router       /doi-settings/{journalId} [post]
func (this *DoiSettingCreateHandler) Handle(ctx echo.Context) error {
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

	return c.JsonResponse(201, "DOI setting created successfully")
}
