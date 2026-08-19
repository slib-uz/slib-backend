package draft

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/draftusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type GetDraftHandler struct {
	uc *draftusecases.GetDraftUseCase
}

// @inject
func NewGetDraftHandler(uc *draftusecases.GetDraftUseCase) *GetDraftHandler {
	return &GetDraftHandler{uc: uc}
}

// Handle
// @Summary Get Draft
// @Description Retrieve a draft by its key
// @Tags draft
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param key query string true "Draft Key"
// @Success 200 {object} map[string]any
// @Router /draft/get [get]
func (this *GetDraftHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	key := c.QueryParam("key")
	if key == "" {
		return c.JSON(400, map[string]any{
			"status":  "400",
			"message": "Key query parameter is required",
		})
	}

	draft, err := this.uc.Execute(key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status": "200",
		"data":   draft,
	})
}
