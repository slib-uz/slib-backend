package auth

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/orcidusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/auth/schema"
)

type OrcidConnectHandler struct {
	uc *orcidusecases.OrcidConnectUseCase
}

// @inject
func NewOrcidConnectHandler(uc *orcidusecases.OrcidConnectUseCase) *OrcidConnectHandler {
	return &OrcidConnectHandler{uc: uc}
}

// Handle godoc
// @Tags orcid
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body schema.OrcidConnectRequest true "Orcid Connect Request"
// @Success 200 {object} response.Response
// @Router /auth/orcid/connect [post]
func (this *OrcidConnectHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	data, err := context.GetBody[schema.OrcidConnectRequest](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User.ID, data.Code, data.RedirectUri); err != nil {
		return err
	}
	return c.JSON(200, map[string]any{
		"status":  200,
		"message": "ORCID account connected successfully",
	})
}
