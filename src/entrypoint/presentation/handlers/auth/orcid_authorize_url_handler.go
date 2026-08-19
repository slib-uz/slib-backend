package auth

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/orcidusecases"
)

type OrcidAuthorizeURLHandler struct {
	uc *orcidusecases.OrcidAuthorizeUriUseCase
}

// @inject
func NewOrcidAuthorizeURLHandler(uc *orcidusecases.OrcidAuthorizeUriUseCase) *OrcidAuthorizeURLHandler {
	return &OrcidAuthorizeURLHandler{uc: uc}
}

// Handle godoc
// @Tags orcid
// @Accept json
// @Produce json
// @Param redirect_uri query string true "Redirect URI"
// @Success 200 {object} response.Response
// @Router /auth/orcid/authorize/url [get]
func (this *OrcidAuthorizeURLHandler) Handle(ctx echo.Context) error {
	authorizeUri := this.uc.Execute(ctx.QueryParam("redirect_uri"))

	return ctx.JSON(200, map[string]any{
		"status": 200,
		"data": map[string]any{
			"redirect_uri":  ctx.QueryParam("redirect_uri"),
			"authorize_uri": authorizeUri,
		},
	})
}
