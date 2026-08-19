package groups

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/entrypoint/presentation/handlers/comment"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
)

type CommentGroup struct {
	commentCreateHandler   *CommentCreateHandler
	articleCommentsHandler *ArticleCommentsListHandler
	commentStatsHandler    *CommentStatsHandler
	commentDeleteHandler   *CommentDeleteHandler
}

// @inject
func NewCommentGroup(commentCreateHandler *CommentCreateHandler, articleCommentsHandler *ArticleCommentsListHandler, commentStatsHandler *CommentStatsHandler, commentDeleteHandler *CommentDeleteHandler) *CommentGroup {
	return &CommentGroup{commentCreateHandler: commentCreateHandler, articleCommentsHandler: articleCommentsHandler, commentStatsHandler: commentStatsHandler, commentDeleteHandler: commentDeleteHandler}
}

func (this *CommentGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/create", this.commentCreateHandler.Handle, permissions.AuthenticatedPermission)
	group.GET("/article/:articleId/list", this.articleCommentsHandler.Handle)
	group.GET("/article/:articleId/stats", this.commentStatsHandler.Handle)
	authGroup := group.Group("")
	authGroup.DELETE("/:id", this.commentDeleteHandler.Handle, permissions.AuthenticatedPermission)
}
