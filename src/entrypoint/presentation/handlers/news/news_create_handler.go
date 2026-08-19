package news

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/newsusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/news/schema"
)

type NewsCreateHandler struct {
	uc *NewsCreateUseCase
}

// @inject
func NewNewsCreateHandler(uc *NewsCreateUseCase) *NewsCreateHandler {
	return &NewsCreateHandler{uc: uc}
}

// Handle NewsCreateHandler
// @Tags news
// @Accept json
// @Produce json
// @Param request body schema.NewsCreateRequest true "Create News Request"
// @Success 201 {object} entity.NewsEntity
// @Router /news/create [post]
func (this *NewsCreateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	request, err := context2.GetBody[schema.NewsCreateRequest](c)
	if err != nil {
		return err
	}

	newsEntity := schema.NewsCreateRequestToEntity(request)

	news, err := this.uc.Execute(c.Request().Context(), newsEntity)
	if err != nil {
		return err
	}

	return c.JsonResponse(201, news)
}
