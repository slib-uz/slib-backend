package article_application

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/article_applications_usecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/article_application/schema"
)

type ApplicationPublishHandler struct {
	uc *ApplicationPublishUseCase
}

// @inject
func NewApplicationPublishHandler(uc *ApplicationPublishUseCase) *ApplicationPublishHandler {
	return &ApplicationPublishHandler{uc: uc}
}

// Handle
// @Tags         article-application
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  schema.ApplicationPublishSchema  true  "Application Publish Schema"
// @Success      200  {object}  map[string]any
// @Router       /article-application/publish [post]
func (this *ApplicationPublishHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	data, err := context.GetBody[schema.ApplicationPublishSchema](ctx)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User, data.ApplicationID, data.FinalFile); err != nil {
		return err
	}

	return c.JsonResponse(200, "Application published successfully")
}
