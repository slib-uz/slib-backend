package upload

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/uploadusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/upload/schema"
)

type PresignedPutUrlHandler struct {
	uc *uploadusecases.UploadFileUseCase
}

// @inject
func NewPresignedPutUrlHandler(uc *uploadusecases.UploadFileUseCase) *PresignedPutUrlHandler {
	return &PresignedPutUrlHandler{uc: uc}
}

// Handle godoc
// @Tags upload
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body schema.PresignedPutUrlRequest true "Upload File Request"
// @Success 200 {object} schema.PresignedPutUrlResponse
// @Failure 400 {object} response.Response "kengaytma ruxsat etilmagan yoki papka noma'lum"
// @Router /upload/presigned/put-url [post]
func (this *PresignedPutUrlHandler) Handle(c echo.Context) error {
	body, err := context.GetBody[schema.PresignedPutUrlRequest](c)
	if err != nil {
		return err
	}

	result, err := this.uc.Execute(c.Request().Context(), body.Folder, body.FileName)
	if err != nil {
		return err
	}

	return c.JSON(200, schema.NewPresignedPutUrlResponse(result.URL, result.Fields, result.ObjectName))
}
