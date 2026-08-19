package publisher

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/publisherusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/publisher/schema"
)

type PublisherCreateHandler struct {
	uc *publisherusecases.PublisherCreateUseCase
}

// @inject
func NewPublisherCreateHandler(uc *publisherusecases.PublisherCreateUseCase) *PublisherCreateHandler {
	return &PublisherCreateHandler{uc: uc}
}

// Handle handles the request to create a new publisher
// @Tags Publisher
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param publisher body schema.PublisherCreateRequest true "Publisher"
// @Success 201 {object} schema.PublisherCreateRequest
// @Router /publisher/create [post]
func (this *PublisherCreateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	data, err := context.GetBody[schema.PublisherCreateRequest](c)
	if err != nil {
		return err
	}

	publisherMap := schema.PublisherResToDto(data)
	if err := this.uc.Execute(publisherMap); err != nil {
		return err
	}
	return c.JsonResponse(201, "Publisher created successfully")
}
