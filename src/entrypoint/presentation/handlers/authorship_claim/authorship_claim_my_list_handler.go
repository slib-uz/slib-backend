package authorship_claim

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/authorshipclaimusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type AuthorshipClaimMyListHandler struct {
	uc *authorshipclaimusecases.ListAuthorshipClaimsUseCase
}

// @inject
func NewAuthorshipClaimMyListHandler(uc *authorshipclaimusecases.ListAuthorshipClaimsUseCase) *AuthorshipClaimMyListHandler {
	return &AuthorshipClaimMyListHandler{uc: uc}
}

// Handle
// @Tags AuthorshipClaim
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page"
// @Param page_size query int false "Page Size"
// @Param status query string false "Status"
// @Success 200 {array} entity.AuthorshipClaimEntity
// @Router /authorship-claims/my [get]
func (this *AuthorshipClaimMyListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	if c.User == nil {
		return echo.ErrUnauthorized
	}

	page, size := context2.GetPagingParams(ctx)
	filters := make(map[string]interface{})

	filters["sender_id"] = c.User.ID

	if status := c.QueryParam("status"); status != "" {
		filters["status"] = status
	}

	result, err := this.uc.Execute(page, size, filters)
	if err != nil {
		return err
	}

	return c.JsonResponse(http.StatusOK, result)
}
