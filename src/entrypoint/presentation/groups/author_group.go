package groups

import (
	"github.com/labstack/echo/v4"
	author2 "slib.uz/src/entrypoint/presentation/handlers/author"
)

//import (
//	container "slib.uz/src/di"
//	"slib.uz/src/presentation/app"
//	"slib.uz/src/presentation/handlers/author"
//)
//
//var AuthorGroup = app.GetInstance().Group("/author")
//
//func init() {
//	AuthorGroup.GET("/find-by-id/:science_id", container.InitAuthorByScienceIDHandler().Handle)
//	AuthorGroup.GET("/profile/:scienceId", container.InitAuthorProfileHandler().Handle)
//	AuthorGroup.GET("/list", container.InitAuthorsListHandler().Handle)
//	AuthorGroup.GET("/top", container.InitTopAuthorsHandler().Handle)
//}

type AuthorGroup struct {
	authorByScienceIDHandler *author2.AuthorByScienceIDHandler
	authorProfileHandler     *author2.AuthorProfileHandler
	authorsListHandler       *author2.AuthorsListHandler
	topAuthorsHandler        *author2.TopAuthorsHandler
}

// @inject
func NewAuthorGroup(authorByScienceIDHandler *author2.AuthorByScienceIDHandler, authorProfileHandler *author2.AuthorProfileHandler, authorsListHandler *author2.AuthorsListHandler, topAuthorsHandler *author2.TopAuthorsHandler) *AuthorGroup {
	return &AuthorGroup{authorByScienceIDHandler: authorByScienceIDHandler, authorProfileHandler: authorProfileHandler, authorsListHandler: authorsListHandler, topAuthorsHandler: topAuthorsHandler}
}

func (this *AuthorGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/find-by-id/:science_id", this.authorByScienceIDHandler.Handle)
	group.GET("/profile/:scienceId", this.authorProfileHandler.Handle)
	group.GET("/list", this.authorsListHandler.Handle)
	group.GET("/top", this.topAuthorsHandler.Handle)
}
