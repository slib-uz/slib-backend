package article_author_affiliation

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/articleauthoraffiliationusecase"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/article_author_affiliation/schema"
)

type ArticleAuthorAffiliationCreateHandler struct {
	uc *ArticleAuthorAffiliationCreateUseCase
}

// @inject
func NewArticleAuthorAffiliationCreateHandler(uc *ArticleAuthorAffiliationCreateUseCase) *ArticleAuthorAffiliationCreateHandler {
	return &ArticleAuthorAffiliationCreateHandler{uc: uc}
}

// Handle ArticleAuthorAffiliationCreateHandler
// @Tags article-author-affiliation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param articleAuthorAffiliation body schema.ArticleAuthorAffiliationCreateRequest true "ArticleAuthorAffiliationCreateRequest"
// @Success 201 {object} response.Response
// @Router /article-author-affiliation/create [post]
func (this *ArticleAuthorAffiliationCreateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	data, err := context2.GetBody[schema.ArticleAuthorAffiliationCreateRequest](c)
	if err != nil {
		return err
	}

	articleAuthorAffiliation, err := this.uc.Execute(data.ToEntity())
	if err != nil {
		return err
	}

	return c.JsonResponse(201, articleAuthorAffiliation)
}
