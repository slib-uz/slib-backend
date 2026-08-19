package upload

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/uploadusecases"
)

type FileUploadHandler struct {
	uc *uploadusecases.UploadTempFileUseCase
}

// @inject
func NewFileUploadHandler(uc *uploadusecases.UploadTempFileUseCase) *FileUploadHandler {
	return &FileUploadHandler{uc: uc}
}

// Handle
// @Accept  multipart/form-data
// @Produce  json
// @Tags upload
// @Security BearerAuth
// @Param file formData file true "ContentFile"
// @Success 200 {object} map[string]string
// @Router /upload/file [post]
func (this *FileUploadHandler) Handle(c echo.Context) error {

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ContentFile must be provided"})
	}

	filePath, err := this.uc.Execute(file)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"file": filePath})
}
