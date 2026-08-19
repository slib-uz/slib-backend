package publisher

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/publisherusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type PublisherDetailHandler struct {
	uc *publisherusecases.PublisherDetailUseCase
}

// @inject
func NewPublisherDetailHandler(uc *publisherusecases.PublisherDetailUseCase) *PublisherDetailHandler {
	return &PublisherDetailHandler{uc: uc}
}

// Handle PublisherDetailHandler handles the request to get publisher details
// @Tags Publisher
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Publisher ID"
// @Success 200 {object} schema.PublisherCreateRequest
// @Router /publisher/detail/{id} [get]
func (this *PublisherDetailHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)
	strId := c.Param("id")
	uintId, err := strconv.ParseUint(strId, 10, 32)

	publisher, err := this.uc.Execute(uint(uintId))
	if err != nil {
		return c.JsonResponse(400, err.Error())
	}
	return c.JsonResponse(200, publisher)
}
