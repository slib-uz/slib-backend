package edition

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/editionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type EditionCreateHandler struct {
	uc *editionusecases.EditionCreateUseCase
}

// @inject
func NewEditionCreateHandler(uc *editionusecases.EditionCreateUseCase) *EditionCreateHandler {
	return &EditionCreateHandler{uc: uc}
}

// Handle godoc
// @Tags         edition
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        journalId  path      int                   true  "Journal ID"
// @Param        request    body      entity.EditionEntity  true  "Edition data"
// @Success      201  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /edition/create [post]
func (this *EditionCreateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	data, err := context.GetBody[entity.EditionEntity](ctx)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(data, c.User); err != nil {
		return err
	}

	return c.JsonResponse(201, "Edition created successfully")
}
