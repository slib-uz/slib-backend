package groups

import (
	"github.com/labstack/echo/v4"
	articles2 "slib.uz/src/entrypoint/presentation/handlers/articles"
	author2 "slib.uz/src/entrypoint/presentation/handlers/author"
	middlewares2 "slib.uz/src/entrypoint/presentation/interceptor/middlewares"
)

// IntegrationGroup — tashqi tizimlar uchun Basic Auth orqali himoyalangan API guruhi.
type IntegrationGroup struct {
	integrationArticleCreateHandler *articles2.IntegrationArticleCreateHandler
	authorByScienceIDHandler        *author2.AuthorByScienceIDHandler

	// middlewares
	basicAuth *middlewares2.ClientBasicAuthDBMiddleware
}

// @inject
func NewIntegrationGroup(
	integrationArticleCreateHandler *articles2.IntegrationArticleCreateHandler,
	authorByScienceIDHandler *author2.AuthorByScienceIDHandler,
	clientBasicAuthDBMiddleware *middlewares2.ClientBasicAuthDBMiddleware,
) *IntegrationGroup {
	return &IntegrationGroup{
		integrationArticleCreateHandler: integrationArticleCreateHandler,
		authorByScienceIDHandler:        authorByScienceIDHandler,
		basicAuth:                       clientBasicAuthDBMiddleware,
	}
}

func (this *IntegrationGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/articles", this.integrationArticleCreateHandler.Handle, this.basicAuth.Wrap)
	group.GET("/authors/find-by-id/:science_id", this.authorByScienceIDHandler.Handle, this.basicAuth.Wrap)
}
