package authorship_claim

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/authorshipclaimusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type AuthorshipClaimListHandler struct {
	uc *authorshipclaimusecases.ListAuthorshipClaimsUseCase
}

// @inject
func NewAuthorshipClaimListHandler(uc *authorshipclaimusecases.ListAuthorshipClaimsUseCase) *AuthorshipClaimListHandler {
	return &AuthorshipClaimListHandler{uc: uc}
}

// Handle
// @Tags AuthorshipClaim
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Param article_id query int false "Article ID"
// @Param journal_id query int false "Journal ID"
// @Param publisher_id query int false "Publisher ID"
// @Param status query string false "Status"
// @Success 200 {array} entity.AuthorshipClaimEntity
// @Router /authorship-claims/list [get]
func (this *AuthorshipClaimListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	page, size := context.GetPagingParams(c)
	filters := make(map[string]interface{})

	if articleID := c.QueryParam("article_id"); articleID != "" {
		filters["article_id"] = articleID
	}
	if journalID := c.QueryParam("journal_id"); journalID != "" {
		filters["journal_id"] = journalID
	}
	if publisherID := c.QueryParam("publisher_id"); publisherID != "" {
		filters["publisher_id"] = publisherID
	}
	if status := c.QueryParam("status"); status != "" {
		filters["status"] = status
	}

	result, err := this.uc.Execute(page, size, filters)
	if err != nil {
		return c.JsonResponse(500, err.Error())
	}

	return c.JsonResponse(200, result)
}
