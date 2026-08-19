package doi

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/doiusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type CrossRefPrefixInfoHandler struct {
	uc *doiusecases.CrossRefPrefixInfoUseCase
}

// @inject
func NewCrossRefPrefixInfoHandler(uc *doiusecases.CrossRefPrefixInfoUseCase) *CrossRefPrefixInfoHandler {
	return &CrossRefPrefixInfoHandler{uc: uc}
}

// Handle
// @Tags         doi-settings
// @Produce      json
// @Security     BearerAuth
// @Param        prefix path string true "DOI Prefix (e.g. 10.53279)"
// @Success      200 {object} map[string]any
// @Failure      400,404,500 {object} response.Response
// @Router       /doi-settings/prefix/{prefix} [get]
func (this *CrossRefPrefixInfoHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)
	prefix := c.Param("prefix")

	info, err := this.uc.Execute(c.Request().Context(), prefix)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, info)
}
