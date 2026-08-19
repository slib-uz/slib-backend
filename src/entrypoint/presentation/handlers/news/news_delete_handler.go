package news

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/newsusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type NewsDeleteHandler struct {
	uc *NewsDeleteUseCase
}

// @inject
func NewNewsDeleteHandler(uc *NewsDeleteUseCase) *NewsDeleteHandler {
	return &NewsDeleteHandler{uc: uc}
}

// Handle NewsDeleteHandler
// @Tags news
// @Accept json
// @Produce json
// @Param newsId path int true "News ID"
// @Success 200 {object} map[string]string
// @Router /news/delete/{newsId} [delete]
func (this *NewsDeleteHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	newsIdValue, err := context2.GetIntPathParam(c, "newsId")
	if err != nil {
		return err
	}
	newsID := uint(newsIdValue)

	err = this.uc.Execute(c.Request().Context(), newsID)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, map[string]string{"message": "News deleted successfully"})
}
