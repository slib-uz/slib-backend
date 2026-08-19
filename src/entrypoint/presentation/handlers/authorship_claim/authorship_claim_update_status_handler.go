package authorship_claim

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/authorshipclaimusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	schema "slib.uz/src/entrypoint/presentation/handlers/authorship_claim/schema"
)

type UpdateAuthorshipClaimStatusHandler struct {
	uc *authorshipclaimusecases.UpdateAuthorshipClaimStatusUseCase
}

// @inject
func NewUpdateAuthorshipClaimStatusHandler(uc *authorshipclaimusecases.UpdateAuthorshipClaimStatusUseCase) *UpdateAuthorshipClaimStatusHandler {
	return &UpdateAuthorshipClaimStatusHandler{uc: uc}
}

// Handle
// @Tags AuthorshipClaim
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Claim ID"
// @Param body body schema.UpdateStatusRequest true "Status Update Body"
// @Success 200 {string} string "OK"
// @Failure 403 {string} string "Forbidden"
// @Failure 400 {string} string "Bad Request"
// @Router /authorship-claims/status/{id} [put]
func (this *UpdateAuthorshipClaimStatusHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	id, err := context2.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}

	body, err := context2.GetBody[schema.UpdateStatusRequest](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User, body.ToEntity(uint(id), c.User.ID)); err != nil {
		return err
	}
	return ctx.JSON(200, map[string]any{
		"message": "Status updated successfully",
	})
}
