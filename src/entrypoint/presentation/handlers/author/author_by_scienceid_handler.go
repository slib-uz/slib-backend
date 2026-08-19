package author

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/authorusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type AuthorByScienceIDHandler struct {
	uc *AuthorByScienceIDUseCase
}

// @inject
func NewAuthorByScienceIDHandler(uc *AuthorByScienceIDUseCase) *AuthorByScienceIDHandler {
	return &AuthorByScienceIDHandler{uc: uc}
}

// Handle handles the request to get author by science ID
// @Tags author
// @Accept json
// @Produce json
// @Param science_id path string true "Science ID"
// @Success 200 {object} entity.AuthorEntity
// @Router /author/find-by-id/{science_id} [get]
func (this *AuthorByScienceIDHandler) Handle(ctx echo.Context) error {

	c := ctx.(*context.Context)

	author, err := this.uc.Execute(c.Param("science_id"))

	if err != nil {
		return err
	}

	return c.JsonResponse(200, author)
}
