package roi

/*

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/roiusecase"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type TestPublishROIHandler struct {
	uc *roiusecase.ROISyncByArticleUseCase
}

// @inject
func NewTestPublishROIHandler(uc *roiusecase.ROISyncByArticleUseCase) *TestPublishROIHandler {
	return &TestPublishROIHandler{uc: uc}
}

// Handle
// @Tags         roi
// @Accept       json
// @Produce      json
// @Security     BasicAuth
// @Param        article_id path int true "Article ID"
// @Success      200  {object}  response.Response
// @Router       /roi/article/test-publish/{article_id} [get]
func (this *TestPublishROIHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)
	articleId, err := context.GetIntPathParam(c, "article_id")
	if err != nil {
		return err
	}

	roi, err := this.uc.Execute(uint(articleId))
	if err != nil {
		return err
	}
	return c.JsonResponse(200, roi)
}

*/
