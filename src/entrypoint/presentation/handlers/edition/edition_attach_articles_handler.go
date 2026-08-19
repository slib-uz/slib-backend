package edition

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/editionusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type editionAttachArticlesRequest struct {
	ArticleIDs []uint `json:"article_ids" validate:"required,min=1"`
}

type EditionAttachArticlesHandler struct {
	uc *editionusecases.EditionAttachArticlesUseCase
}

// @inject
func NewEditionAttachArticlesHandler(uc *editionusecases.EditionAttachArticlesUseCase) *EditionAttachArticlesHandler {
	return &EditionAttachArticlesHandler{uc: uc}
}

// Handle godoc
// @Tags         edition
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        editionId  path      int                             true  "Edition ID"
// @Param        request    body      editionAttachArticlesRequest    true  "Article IDs"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /edition/{editionId}/attach-articles [post]
func (this *EditionAttachArticlesHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	editionIdValue, err := context.GetIntPathParam(ctx, "editionId")
	if err != nil {
		return err
	}
	editionID := uint(editionIdValue)

	body, err := context.GetBody[editionAttachArticlesRequest](ctx)
	if err != nil {
		return err
	}

	affected, err := this.uc.Execute(ctx.Request().Context(), editionID, body.ArticleIDs, c.User)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, map[string]int64{"attached_count": affected})
}
