package roi

/*

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/roiusecase"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type ArticleROIHandler struct {
	uc *roiusecase.ArticleROISearchUseCase
}

// @inject
func NewArticleROIHandler(uc *roiusecase.ArticleROISearchUseCase) *ArticleROIHandler {
	return &ArticleROIHandler{uc: uc}
}

// Handle
// @Tags         roi
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        app_id path int true "Application ID"
// @Success      200  {object}  response.Response
// @Router       /roi/article/search/{app_id} [get]
func (this *ArticleROIHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)
	appId, err := context.GetIntPathParam(c, "app_id")
	if err != nil {
		return err
	}

	result, err := this.uc.Execute(uint(appId))
	if err != nil {
		return response.NewOptionalResponse(200, response.RoiNotFound, err.Error(), "")
	}
	return response.NewOptionalResponse(200, response.RoiFound, result, "")
}

*/
