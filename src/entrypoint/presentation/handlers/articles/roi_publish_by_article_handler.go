package articles

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/roiusecase"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type ROIPublishByArticleHandler struct {
	uc *roiusecase.ROIPublishByArticleUseCase
}

// @inject
func NewROIPublishByArticleHandler(uc *roiusecase.ROIPublishByArticleUseCase) *ROIPublishByArticleHandler {
	return &ROIPublishByArticleHandler{uc: uc}
}

// Handle godoc
// @Tags         article
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        articleId  path  int  true  "Article ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /articles/roi-publish/{articleId} [post]
func (this *ROIPublishByArticleHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	articleId, err := context.GetIntPathParam(ctx, "articleId")
	if err != nil {
		return err
	}

	roi, err := this.uc.Execute(uint(articleId))
	if err != nil {
		return err
	}

	return c.JsonResponse(200, roi)
}
