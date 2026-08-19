package groups

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/entrypoint/presentation/handlers/article_author_affiliation"
)

type ArticleAuthorAffiliationGroup struct {
	ArticleAuthorAffiliationCreateHandler *ArticleAuthorAffiliationCreateHandler
}

// @inject
func NewArticleAuthorAffiliationGroup(articleAuthorAffiliationCreateHandler *ArticleAuthorAffiliationCreateHandler) *ArticleAuthorAffiliationGroup {
	return &ArticleAuthorAffiliationGroup{ArticleAuthorAffiliationCreateHandler: articleAuthorAffiliationCreateHandler}
}

func (this *ArticleAuthorAffiliationGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/create", this.ArticleAuthorAffiliationCreateHandler.Handle)
}
