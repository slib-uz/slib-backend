package authorship_claim

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/authorshipclaimusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type AuthorshipClaimDetailHandler struct {
	uc *authorshipclaimusecases.GetAuthorshipClaimDetailUseCase
}

// @inject
func NewAuthorshipClaimDetailHandler(uc *authorshipclaimusecases.GetAuthorshipClaimDetailUseCase) *AuthorshipClaimDetailHandler {
	return &AuthorshipClaimDetailHandler{uc: uc}
}

// Handle
// @Tags AuthorshipClaim
// @Accept json
// @Produce json
// @Param id path int true "Authorship Claim ID"
// @Success 200 {object} entity.AuthorshipClaimEntity
// @Failure 404 {object} response.Response
// @Router /authorship-claims/detail/{id} [get]
func (this *AuthorshipClaimDetailHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JsonResponse(400, "Invalid ID")
	}

	claim, err := this.uc.Execute(uint(id))
	if err != nil {
		// Basic error handling, improvement: custom errors for Not Found
		return c.JsonResponse(404, "Claim not found")
	}
	if claim == nil {
		return c.JsonResponse(404, "Claim not found")
	}

	return c.JsonResponse(200, claim)
}
