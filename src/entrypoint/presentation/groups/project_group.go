package groups

import (
	"github.com/labstack/echo/v4"
	project2 "slib.uz/src/entrypoint/presentation/handlers/project"
)

type ProjectGroup struct {
	projectsListHandler  *project2.ProjectsListHandler
	projectDetailHandler *project2.ProjectDetailHandler
}

// @inject
func NewProjectGroup(projectListHandler *project2.ProjectsListHandler, projectDetailHandler *project2.ProjectDetailHandler) *ProjectGroup {
	return &ProjectGroup{
		projectsListHandler:  projectListHandler,
		projectDetailHandler: projectDetailHandler,
	}
}

func (this *ProjectGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/list", this.projectsListHandler.Handle)
	group.GET("/detail/:id", this.projectDetailHandler.Handle)
}
