package edition

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/editionusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type EditionDetachArticlesHandler struct {
	uc *editionusecases.EditionDetachArticlesUseCase
}

// @inject
func NewEditionDetachArticlesHandler(uc *editionusecases.EditionDetachArticlesUseCase) *EditionDetachArticlesHandler {
	return &EditionDetachArticlesHandler{uc: uc}
}

type editionDetachArticlesRequest struct {
	ArticleIDs []uint `json:"article_ids" validate:"required,min=1"`
}

// Handle godoc
// @Tags         edition
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        editionId  path      int                             true  "Edition ID"
// @Param        request    body      editionDetachArticlesRequest    true  "Article IDs"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /edition/{editionId}/detach-articles [post]
func (this *EditionDetachArticlesHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	editionIdValue, err := context.GetIntPathParam(ctx, "editionId")
	if err != nil {
		return err
	}
	editionID := uint(editionIdValue)

	body, err := context.GetBody[editionDetachArticlesRequest](ctx)
	if err != nil {
		return err
	}

	affected, err := this.uc.Execute(ctx.Request().Context(), editionID, body.ArticleIDs, c.User)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, map[string]int64{"detached_count": affected})
}
