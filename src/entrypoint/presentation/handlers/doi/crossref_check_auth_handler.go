package doi

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/doiusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/doi/schema"
)

type CrossRefCheckAuthHandler struct {
	uc *doiusecases.CrossRefCheckAuthUseCase
}

// @inject
func NewCrossRefCheckAuthHandler(uc *doiusecases.CrossRefCheckAuthUseCase) *CrossRefCheckAuthHandler {
	return &CrossRefCheckAuthHandler{uc: uc}
}

// Handle
// @Tags         doi-settings
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body schema.CrossRefCheckAuthRequest true "CrossRef credentials"
// @Success      200 {object} map[string]any
// @Failure      400,401,403,404,500 {object} response.Response
// @Router       /doi-settings/check-auth [post]
func (this *CrossRefCheckAuthHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	body, err := context.GetBody[schema.CrossRefCheckAuthRequest](c)
	if err != nil {
		return err
	}

	authenticated, err := this.uc.Execute(c.Request().Context(), body.Username, body.Password)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, map[string]bool{"authenticated": authenticated})
}
