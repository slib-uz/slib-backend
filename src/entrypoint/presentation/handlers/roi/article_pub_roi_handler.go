package roi

/*

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/roiusecase"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/article_application/schema"
)

type ArticlePubRoiHandler struct {
	uc *roiusecase.ArticlePushROIUseCase
}

// @inject
func NewArticlePubRoiHandler(uc *roiusecase.ArticlePushROIUseCase) *ArticlePubRoiHandler {
	return &ArticlePubRoiHandler{uc: uc}
}

// Handle
// @Tags         roi
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  schema.ApplicationPublishSchema  true  "Application Publish Schema"
// @Success      200  {object}  response.Response
// @Router       /roi/article/publish [post]
func (this *ArticlePubRoiHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	data, err := context.GetBody[schema.ApplicationPublishSchema](ctx)
	if err != nil {
		return err
	}

	roi, err := this.uc.Execute(c.User, data.ApplicationID, data.FinalFile)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, roi)
}

*/
