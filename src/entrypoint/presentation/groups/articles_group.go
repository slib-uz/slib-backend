package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/domain/entity/enum"
	articles2 "slib.uz/src/entrypoint/presentation/handlers/articles"
	middlewares2 "slib.uz/src/entrypoint/presentation/interceptor/middlewares"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
)

type ArticlesGroup struct {
	articleDetailHandler          *articles2.ArticleDetailHandler
	userArticleListHandler        *articles2.UserArticlesListHandler
	journalArticlesListHandler    *articles2.JournalArticlesListHandler
	articlesListHandler           *articles2.ArticlesListHandler
	articleDownloadHandler        *articles2.ArticleDownloadHandler
	articleAndJournalCountHandler *articles2.ArticleAndJournalCountHandler
	topArticlesHandler            *articles2.TopArticlesHandler
	statisticsByStudyFieldHandler *articles2.ArticleStatisticsByStudyFieldHandler
	statisticsByMonthHandler      *articles2.ArticleStatisticsByMonthHandler
	statisticsByYearHandler         *articles2.ArticleStatisticsByYearHandler
	articleUpdateHandler          *articles2.ArticleUpdateHandler
	articleListOnlyNameHandler    *articles2.ArticleListOnlyNameHandler
	directArticleCreateHandler    *articles2.DirectArticleCreateHandler
	roiBackfillHandler            *articles2.ROIBackfillHandler
	roiPublishByArticleHandler    *articles2.ROIPublishByArticleHandler

	// middlewares
	anonymAuthMiddleware *middlewares2.JwAnonymAuthMiddleware
	basicAuthMiddleware  *middlewares2.DeveloperUserBasicAuthMiddleware
}

// @inject
func NewArticlesGroup(
	articleDetailHandler *articles2.ArticleDetailHandler,
	userArticleListHandler *articles2.UserArticlesListHandler,
	journalArticlesListHandler *articles2.JournalArticlesListHandler,
	articlesListHandler *articles2.ArticlesListHandler,
	anonymAuthMiddleware *middlewares2.JwAnonymAuthMiddleware,
	articleDownloadHandler *articles2.ArticleDownloadHandler,
	basicAuth *middlewares2.DeveloperUserBasicAuthMiddleware,
	articleAndJournalCountHandler *articles2.ArticleAndJournalCountHandler,
	topArticlesHandler *articles2.TopArticlesHandler,
	statisticsByStudyFieldHandler *articles2.ArticleStatisticsByStudyFieldHandler,
	statisticsByMonthHandler *articles2.ArticleStatisticsByMonthHandler,
	statisticsByYearHandler *articles2.ArticleStatisticsByYearHandler,
	articleUpdateHandler *articles2.ArticleUpdateHandler,
	articleListOnlyNameHandler *articles2.ArticleListOnlyNameHandler,
	directArticleCreateHandler *articles2.DirectArticleCreateHandler,
	roiBackfillHandler *articles2.ROIBackfillHandler,
	roiPublishByArticleHandler *articles2.ROIPublishByArticleHandler,
) *ArticlesGroup {
	return &ArticlesGroup{
		articleDetailHandler:          articleDetailHandler,
		userArticleListHandler:        userArticleListHandler,
		journalArticlesListHandler:    journalArticlesListHandler,
		articlesListHandler:           articlesListHandler,
		anonymAuthMiddleware:          anonymAuthMiddleware,
		articleDownloadHandler:        articleDownloadHandler,
		basicAuthMiddleware:           basicAuth,
		articleAndJournalCountHandler: articleAndJournalCountHandler,
		topArticlesHandler:            topArticlesHandler,
		statisticsByStudyFieldHandler: statisticsByStudyFieldHandler,
		statisticsByMonthHandler:      statisticsByMonthHandler,
		statisticsByYearHandler:         statisticsByYearHandler,
		articleUpdateHandler:          articleUpdateHandler,
		articleListOnlyNameHandler:    articleListOnlyNameHandler,
		directArticleCreateHandler:    directArticleCreateHandler,
		roiBackfillHandler:            roiBackfillHandler,
		roiPublishByArticleHandler:    roiPublishByArticleHandler,
	}
}

func (this *ArticlesGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/detail/:articleId", this.articleDetailHandler.Handle, this.anonymAuthMiddleware.Call)
	group.GET("/author/:scienceId", this.userArticleListHandler.Handle)
	group.GET("/journal/:journalId", this.journalArticlesListHandler.Handle)
	group.GET("/list", this.articlesListHandler.Handle)
	group.GET("/:articleId/download", this.articleDownloadHandler.Handle, this.basicAuthMiddleware.Wrap)
	group.GET("/top", this.topArticlesHandler.Handle)
	group.GET("/stats", this.articleAndJournalCountHandler.Handle)
	group.GET("/statistics/by-study-field", this.statisticsByStudyFieldHandler.Handle)
	group.GET("/statistics/by-month", this.statisticsByMonthHandler.Handle)
	group.GET("/statistics/by-year", this.statisticsByYearHandler.Handle)
	group.PUT("/update/:articleId", this.articleUpdateHandler.Handle, permissions.AuthenticatedPermission)
	group.GET("/list-only-name", this.articleListOnlyNameHandler.Handle)
	group.POST("/direct-create/:journalId", this.directArticleCreateHandler.Handle, permissions.AuthenticatedPermission)
	group.POST("/roi-backfill", this.roiBackfillHandler.Handle, this.basicAuthMiddleware.Wrap)
	group.POST("/roi-publish/:articleId", permissions.RolePermission(this.roiPublishByArticleHandler.Handle, enum.RoleAdmin))
}
