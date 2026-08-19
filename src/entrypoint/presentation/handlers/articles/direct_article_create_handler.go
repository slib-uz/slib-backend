package articles

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/articleusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/articles/schema"
)

type DirectArticleCreateHandler struct {
	uc *DirectArticleCreateUseCase
}

// @inject
func NewDirectArticleCreateHandler(uc *DirectArticleCreateUseCase) *DirectArticleCreateHandler {
	return &DirectArticleCreateHandler{uc: uc}
}

// Handle godoc
// @Tags         article
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        journalId  path  int  true  "Journal ID"
// @Param        request  body  schema.ArticleCreateSchema  true  "Article create data"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /articles/direct-create/{journalId} [post]
func (this *DirectArticleCreateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	body, err := context.GetBody[schema.ArticleCreateSchema](c)
	if err != nil {
		return err
	}

	journalId, err := context.GetIntPathParam(ctx, "journalId")
	if err != nil {
		return err
	}

	if err := this.uc.Execute(body.ToEntity(), c.User, uint(journalId)); err != nil {
		return err
	}

	return c.JsonResponse(200, "Article created successfully")
}

// test