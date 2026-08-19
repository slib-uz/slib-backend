package news

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/newsusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/news/schema"
)

type NewsUpdateHandler struct {
	uc *NewsUpdateUseCase
}

// @inject
func NewNewsUpdateHandler(uc *NewsUpdateUseCase) *NewsUpdateHandler {
	return &NewsUpdateHandler{uc: uc}
}

// Handle NewsUpdateHandler
// @Tags news
// @Accept json
// @Produce json
// @Param newsId path int true "News ID"
// @Param request body schema.NewsUpdateRequest true "Update News Request"
// @Success 200 {object} map[string]string
// @Router /news/update/{newsId} [put]
func (this *NewsUpdateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	newsIdValue, err := context2.GetIntPathParam(c, "newsId")
	if err != nil {
		return err
	}
	newsID := uint(newsIdValue)

	request, err := context2.GetBody[schema.NewsUpdateRequest](c)
	if err != nil {
		return err
	}

	newsEntity := schema.NewsUpdateRequestToEntity(request)

	if err := this.uc.Execute(c.Request().Context(), newsID, newsEntity); err != nil {
		return err
	}

	return c.JsonResponse(200, map[string]string{"message": "News updated successfully"})
}
