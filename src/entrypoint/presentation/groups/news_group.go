package groups

import (
	"github.com/labstack/echo/v4"
	news2 "slib.uz/src/entrypoint/presentation/handlers/news"
	"slib.uz/src/entrypoint/presentation/interceptor/middlewares"
)

type NewsGroup struct {
	newsListHandler     *news2.NewsListHandler
	newsRetrieveHandler *news2.NewsRetrieveHandler
	newsCreateHandler   *news2.NewsCreateHandler
	newsUpdateHandler   *news2.NewsUpdateHandler
	newsDeleteHandler   *news2.NewsDeleteHandler

	// middlewares
	anonymAuthMiddleware *middlewares.JwAnonymAuthMiddleware
}

// @inject
func NewNewsGroup(newsListHandler *news2.NewsListHandler, newsRetrieveHandler *news2.NewsRetrieveHandler, newsCreateHandler *news2.NewsCreateHandler, newsUpdateHandler *news2.NewsUpdateHandler, newsDeleteHandler *news2.NewsDeleteHandler, anonymAuthMiddleware *middlewares.JwAnonymAuthMiddleware) *NewsGroup {
	return &NewsGroup{newsListHandler: newsListHandler, newsRetrieveHandler: newsRetrieveHandler, newsCreateHandler: newsCreateHandler, newsUpdateHandler: newsUpdateHandler, newsDeleteHandler: newsDeleteHandler, anonymAuthMiddleware: anonymAuthMiddleware}
}

func (this NewsGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/list", this.newsListHandler.Handle)
	group.GET("/retrieve/:newsId", this.newsRetrieveHandler.Handle, this.anonymAuthMiddleware.Call)
	group.POST("/create", this.newsCreateHandler.Handle)
	group.PUT("/update/:newsId", this.newsUpdateHandler.Handle)
	group.DELETE("/delete/:newsId", this.newsDeleteHandler.Handle)
}
