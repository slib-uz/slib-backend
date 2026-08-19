package upload

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/uploadusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type PresignedGetUrlHandler struct {
	uc *PresignedURLUseCase
}

// @inject
func NewPresignedGetUrlHandler(uc *PresignedURLUseCase) *PresignedGetUrlHandler {
	return &PresignedGetUrlHandler{uc: uc}
}

// Handle
// @Tags upload
// @Accept json
// @Produce json
// @Param articleId path int true "Article ID"
// @Success 200 {object} response.Response
// @Router /upload/presigned-url/article/{articleId} [get]
func (this *PresignedGetUrlHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	articleId, err := context2.GetIntPathParam(c, "articleId")
	if err != nil {
		return err
	}

	url, err := this.uc.Execute(c.User, uint(articleId))
	if err != nil {
		return err
	}
	return c.JsonResponse(200, map[string]string{"url": url})
}
