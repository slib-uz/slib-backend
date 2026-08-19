package edition

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/editionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type EditionUpdateHandler struct {
	uc *editionusecases.EditionUpdateUseCase
}

// @inject
func NewEditionUpdateHandler(uc *editionusecases.EditionUpdateUseCase) *EditionUpdateHandler {
	return &EditionUpdateHandler{uc: uc}
}

// Handle godoc
// @Tags         edition
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        editionId  path      int                   true  "Edition ID"
// @Param        request    body      entity.EditionEntity  true  "Edition data"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /edition/{editionId}/update [put]
func (this *EditionUpdateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	editionIdValue, err := context.GetIntPathParam(ctx, "editionId")
	if err != nil {
		return err
	}
	editionID := uint(editionIdValue)

	data, err := context.GetBody[entity.EditionEntity](ctx)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(ctx.Request().Context(), editionID, data, c.User); err != nil {
		return err
	}

	return c.JsonResponse(200, "Edition updated successfully")
}
