package edition

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/editionusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type EditionDetailHandler struct {
	uc *editionusecases.EditionDetailUseCase
}

// @inject
func NewEditionDetailHandler(uc *editionusecases.EditionDetailUseCase) *EditionDetailHandler {
	return &EditionDetailHandler{uc: uc}
}

// Handle godoc
// @Tags         edition
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Edition ID"
// @Success      200  {object}  entity.EditionEntity
// @Failure      400  {object}  response.Response
// @Router       /edition/{id}/detail [get]
func (this *EditionDetailHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	idValue, err := context.GetIntPathParam(ctx, "id")
	if err != nil {
		return err
	}
	id := uint(idValue)

	result, err := this.uc.Execute(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
