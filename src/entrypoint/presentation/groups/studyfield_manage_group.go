package groups

import (
	"github.com/labstack/echo/v4"
	studyfield2 "slib.uz/src/entrypoint/presentation/handlers/studyfield"
)

//var StudyFieldManageGroup = app.GetInstance().Group(
//	"/studyfield/manage",
//	container.InitJwtAuthMiddleware().Call,
//	permissions.AdminPermission,
//)

//func init() {
//	StudyFieldManageGroup.POST("/create", container.InitStudyFieldCreateHandler().Handle)
//	StudyFieldManageGroup.DELETE("/delete/:id", container.InitStudyFieldDeleteHandler().Handle)
//	StudyFieldManageGroup.PATCH("/update", container.InitStudyFieldUpdateHandler().Handle)
//}

type StudyFieldManageGroup struct {
	createHandler *studyfield2.StudyFieldCreateHandler
	deleteHandler *studyfield2.StudyFieldDeleteHandler
	updateHandler *studyfield2.StudyFieldUpdateHandler
}

// @inject
func NewStudyFieldManageGroup(createHandler *studyfield2.StudyFieldCreateHandler, deleteHandler *studyfield2.StudyFieldDeleteHandler, updateHandler *studyfield2.StudyFieldUpdateHandler) *StudyFieldManageGroup {
	return &StudyFieldManageGroup{createHandler: createHandler, deleteHandler: deleteHandler, updateHandler: updateHandler}
}

func (this *StudyFieldManageGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/create", this.createHandler.Handle)
	group.DELETE("/delete/:id", this.deleteHandler.Handle)
	group.PATCH("/update", this.updateHandler.Handle)
}
