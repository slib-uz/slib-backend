package authorship_claim

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/authorshipclaimusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/authorship_claim/schema"
)

type CreateAuthorshipClaimHandler struct {
	uc *authorshipclaimusecases.CreateAuthorshipClaimUseCase
}

// @inject
func NewCreateAuthorshipClaimHandler(uc *authorshipclaimusecases.CreateAuthorshipClaimUseCase) *CreateAuthorshipClaimHandler {
	return &CreateAuthorshipClaimHandler{uc: uc}
}

// Handle
// @Tags         AuthorshipClaim
// @Accept       json
// @Produce      json
// @Param        request body schema.CreateAuthorshipClaimRequest true "Request body"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /authorship-claims/create [post]
func (this *CreateAuthorshipClaimHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	data, err := context2.GetBody[schema.CreateAuthorshipClaimRequest](ctx)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User, data.ToEntity(c.User.ID)); err != nil {
		return err
	}

	return c.JsonResponse(200, "success")
}
