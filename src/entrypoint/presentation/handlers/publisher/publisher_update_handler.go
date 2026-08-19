package publisher

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/publisherusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	schema2 "slib.uz/src/entrypoint/presentation/handlers/publisher/schema"
)

type PublisherUpdateHandler struct {
	uc *publisherusecases.PublisherUpdateUseCase
}

// @inject
func NewPublisherUpdateHandler(uc *publisherusecases.PublisherUpdateUseCase) *PublisherUpdateHandler {
	return &PublisherUpdateHandler{uc: uc}
}

// Handle PublisherUpdateHandler
// @Tags Publisher
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Publisher ID"
// @Param publisher body schema.PublisherCreateRequest true "Publisher"
// @Success 200 {object} schema.PublisherCreateRequest
// @Router /publisher/update/{id} [put]
func (this *PublisherUpdateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	strId := c.Param("id")
	uintId, err := strconv.ParseUint(strId, 10, 32)
	if err != nil {
		return c.JsonResponse(http.StatusBadRequest, "Invalid ID")
	}

	data, err := context2.GetBody[schema2.PublisherCreateRequest](c)
	if err != nil {
		return err
	}

	publisherMap := schema2.PublisherResToDto(data)

	if err := this.uc.Execute(uint(uintId), publisherMap); err != nil {
		return err
	}
	return c.JsonResponse(http.StatusOK, "Publisher updated successfully")
}
