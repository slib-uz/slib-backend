package draft

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/draftusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/draft/schema"
)

type DraftSaveHandler struct {
	uc *draftusecases.DraftSaveUseCase
}

// @inject
func NewDraftSaveHandler(uc *draftusecases.DraftSaveUseCase) *DraftSaveHandler {
	return &DraftSaveHandler{uc: uc}
}

// Handle
// @Summary Save Draft
// @Description Save a draft
// @Tags draft
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param draft body schema.DraftSaveReqeust true "Draft Save Request"
// @Success 200 {object} map[string]any
// @Router /draft/save [post]
func (this *DraftSaveHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	body, err := context2.GetBody[schema.DraftSaveReqeust](ctx)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(body.ToEntity(c.User.ID)); err != nil {
		return err
	}

	return c.JSON(200, map[string]any{
		"status":  "200",
		"message": "Draft saved successfully",
	})
}
