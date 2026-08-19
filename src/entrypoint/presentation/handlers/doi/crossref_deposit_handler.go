package doi

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/doiusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type CrossRefDepositHandler struct {
	uc *doiusecases.CrossRefDepositUseCase
}

// @inject
func NewCrossRefDepositHandler(uc *doiusecases.CrossRefDepositUseCase) *CrossRefDepositHandler {
	return &CrossRefDepositHandler{uc: uc}
}

// Handle
// @Tags         doi-settings
// @Security     BearerAuth
// @Produce      json
// @Param        articleId path int true "Article ID"
// @Success      200 {object} map[string]any
// @Failure      400,401,403,404,500 {object} response.Response
// @Router       /doi-settings/deposit/{articleId} [post]
func (this *CrossRefDepositHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	articleID, err := context.GetIntPathParam(c, "articleId")
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.Request().Context(), uint(articleID)); err != nil {
		return err
	}

	return c.JsonResponse(200, "CrossRef deposit submitted successfully")
}
