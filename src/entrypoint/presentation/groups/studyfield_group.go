package groups

import (
	"github.com/labstack/echo/v4"
	studyfield2 "slib.uz/src/entrypoint/presentation/handlers/studyfield"
)

//
//var StudyFieldGroup = app.GetInstance().Group("/studyfield")
//
//var StudyFieldManageGroup = app.GetInstance().Group(
//	"/studyfield/manage",
//	container.InitJwtAuthMiddleware().Call,
//	permissions.AdminPermission,
//)

//func init() {
//	StudyFieldGroup.GET("/list", container.InitStudyFieldListHandler().Handle)
//	StudyFieldGroup.GET("/paging", container.InitStudyFieldPagingHandler().Handle)
//
//	// manage
//	StudyFieldManageGroup.POST("/create", container.InitStudyFieldCreateHandler().Handle)
//	StudyFieldManageGroup.DELETE("/delete/:id", container.InitStudyFieldDeleteHandler().Handle)
//	StudyFieldManageGroup.PATCH("/update", container.InitStudyFieldUpdateHandler().Handle)
//}

type StudyFieldGroup struct {
	listHandler   *studyfield2.StudyFieldListHandler
	pagingHandler *studyfield2.StudyFieldPagingHandler
}

// @inject
func NewStudyFieldGroup(listHandler *studyfield2.StudyFieldListHandler, pagingHandler *studyfield2.StudyFieldPagingHandler) *StudyFieldGroup {
	return &StudyFieldGroup{listHandler: listHandler, pagingHandler: pagingHandler}
}

func (this *StudyFieldGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/list", this.listHandler.Handle)
	group.GET("/paging", this.pagingHandler.Handle)
}
