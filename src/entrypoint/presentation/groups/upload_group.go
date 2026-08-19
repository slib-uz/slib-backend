package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/upload"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
)

type UploadGroup struct {
	fileUploadHandler      *upload.FileUploadHandler
	presignedGetUrlHandler *upload.PresignedGetUrlHandler
	presignedPutUrlHandler *upload.PresignedPutUrlHandler
}

// @inject
func NewUploadGroup(fileUploadHandler *upload.FileUploadHandler, presignedGetUrlHandler *upload.PresignedGetUrlHandler, presignedPutUrlHandler *upload.PresignedPutUrlHandler) *UploadGroup {
	return &UploadGroup{fileUploadHandler: fileUploadHandler, presignedGetUrlHandler: presignedGetUrlHandler, presignedPutUrlHandler: presignedPutUrlHandler}
}

func (this *UploadGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/file", this.fileUploadHandler.Handle, permissions.AuthenticatedPermission)
	group.GET("/presigned-url/article/:articleId", this.presignedGetUrlHandler.Handle)
	group.POST("/presigned/put-url", this.presignedPutUrlHandler.Handle, permissions.AuthenticatedPermission)
}
