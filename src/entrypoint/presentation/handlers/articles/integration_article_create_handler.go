package articles

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/articleusecases"
	"slib.uz/src/core/application/response"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/articles/schema"
)

type IntegrationArticleCreateHandler struct {
	uc *IntegrationArticleCreateUseCase
}

// @inject
func NewIntegrationArticleCreateHandler(
	uc *IntegrationArticleCreateUseCase,
) *IntegrationArticleCreateHandler {
	return &IntegrationArticleCreateHandler{
		uc: uc,
	}
}

// Handle godoc
// @Tags         integration
// @Accept       json
// @Produce      json
// @Security BasicAuth
// @Param        request  body  schema.IntegrationArticleCreateSchema  true  "Article create data"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /integration/articles [post]
func (this *IntegrationArticleCreateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	body, err := context.GetBody[schema.IntegrationArticleCreateSchema](c)
	if err != nil {
		return err
	}

	client := c.GetClient()
	if client == nil || client.JournalID == nil {
		return c.JSON(response.UnauthorizedError.Status, response.UnauthorizedError)
	}

	entity := body.ToEntity()

	if err := this.uc.Execute(entity, *client.JournalID); err != nil {
		return err
	}

	return c.JsonResponse(200, "Article created successfully")
}
