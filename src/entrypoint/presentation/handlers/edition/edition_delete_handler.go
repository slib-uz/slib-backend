package edition

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/editionusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type EditionDeleteHandler struct {
	uc *editionusecases.EditionDeleteUseCase
}

// @inject
func NewEditionDeleteHandler(uc *editionusecases.EditionDeleteUseCase) *EditionDeleteHandler {
	return &EditionDeleteHandler{uc: uc}
}

// Handle godoc
// @Tags         edition
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        editionId  path      int  true  "Edition ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /edition/{editionId}/delete [delete]
func (this *EditionDeleteHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	editionIdValue, err := context.GetIntPathParam(ctx, "editionId")
	if err != nil {
		return err
	}
	editionID := uint(editionIdValue)

	err = this.uc.Execute(ctx.Request().Context(), editionID, c.User)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, map[string]string{"message": "deleted"})
}
