package author

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/authorusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type AuthorProfileHandler struct {
	uc *AuthorProfileUseCase
}

// @inject
func NewAuthorProfileHandler(uc *AuthorProfileUseCase) *AuthorProfileHandler {
	return &AuthorProfileHandler{uc: uc}
}

// Handle
// @Tags author
// @Accept  json
// @Produce  json
// @Param scienceId path string true "Science ID"
// @Success 200 {object} entity.UserProfileEntity
// @Router /author/profile/{scienceId} [get]
func (this *AuthorProfileHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	profile, err := this.uc.Execute(c.Param("scienceId"))
	if err != nil {
		return err
	}

	return c.JsonResponse(200, profile)
}
