package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/authorship_claim"
)

type AuthorshipClaimGroup struct {
	createHandler       *authorship_claim.CreateAuthorshipClaimHandler
	listHandler         *authorship_claim.AuthorshipClaimListHandler
	detailHandler       *authorship_claim.AuthorshipClaimDetailHandler
	myListHandler       *authorship_claim.AuthorshipClaimMyListHandler
	updateStatusHandler *authorship_claim.UpdateAuthorshipClaimStatusHandler
}

// @inject
func NewAuthorshipClaimGroup(
	createHandler *authorship_claim.CreateAuthorshipClaimHandler,
	listHandler *authorship_claim.AuthorshipClaimListHandler,
	detailHandler *authorship_claim.AuthorshipClaimDetailHandler,
	myListHandler *authorship_claim.AuthorshipClaimMyListHandler,
	updateStatusHandler *authorship_claim.UpdateAuthorshipClaimStatusHandler,
) *AuthorshipClaimGroup {
	return &AuthorshipClaimGroup{
		createHandler:       createHandler,
		listHandler:         listHandler,
		detailHandler:       detailHandler,
		myListHandler:       myListHandler,
		updateStatusHandler: updateStatusHandler,
	}
}

func (this *AuthorshipClaimGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/create", this.createHandler.Handle)
	group.GET("/list", this.listHandler.Handle)
	group.GET("/my", this.myListHandler.Handle)
	group.GET("/detail/:id", this.detailHandler.Handle)
	group.PUT("/status/:id", this.updateStatusHandler.Handle)
}
