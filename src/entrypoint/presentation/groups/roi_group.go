package groups

/*

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/roi"
	"slib.uz/src/entrypoint/presentation/interceptor/middlewares"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
)

type RoiGroup struct {
	search                           *roi.ArticleROIHandler
	publish                          *roi.ArticlePubRoiHandler
	test                             *roi.TestPublishROIHandler
	developerUserBasicAuthMiddleware *middlewares.DeveloperUserBasicAuthMiddleware
}

// @inject
func NewRoiGroup(
	search *roi.ArticleROIHandler,
	publish *roi.ArticlePubRoiHandler,
	test *roi.TestPublishROIHandler,
	developerUserBasicAuthMiddleware *middlewares.DeveloperUserBasicAuthMiddleware,
) *RoiGroup {
	return &RoiGroup{
		search:                           search,
		publish:                          publish,
		test:                             test,
		developerUserBasicAuthMiddleware: developerUserBasicAuthMiddleware,
	}
}

func (this *RoiGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/article/search/:app_id", this.search.Handle, permissions.AuthenticatedPermission)
	group.GET("/article/test-publish/:article_id", this.test.Handle, this.developerUserBasicAuthMiddleware.Wrap)
	group.POST("/article/publish", this.publish.Handle, permissions.AuthenticatedPermission)
}

*/