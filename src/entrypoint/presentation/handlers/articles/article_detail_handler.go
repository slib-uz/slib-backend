package articles

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	. "slib.uz/src/core/application/usecase/articleusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type ArticleDetailHandler struct {
	uc *ArticleDetailUseCase
}

// @inject
func NewArticleDetailHandler(uc *ArticleDetailUseCase) *ArticleDetailHandler {
	return &ArticleDetailHandler{uc: uc}
}

// Handle
// @Tags         article
// @Accept       json
// @Produce      json
// @Param        X-Anonymous-Token  header  string  false  "Anon JWT token"
// @Param        articleId  path      int  true  "Article ID"
// @Success      200  {object} entity.ArticleInputEntity
// @Failure      400  {object}  response.Response
// @Router       /articles/detail/{articleId} [get]
func (this *ArticleDetailHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	userKey, err := this.userKey(c)
	if err != nil {
		return err
	}

	articleId, err := context2.GetIntPathParam(c, "articleId")
	if err != nil {
		return err
	}

	article, err := this.uc.Execute(uint(articleId), userKey)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, article)
}

func (this *ArticleDetailHandler) userKey(c *context2.Context) (string, error) {
	if c.User == nil {
		if c.AnonymousID == "" {
			return "", response.NewFailResponse(400, "Anonymous token is required")
		}
		return c.AnonymousID, nil
	}
	return strconv.Itoa(int(c.User.ID)), nil
}
