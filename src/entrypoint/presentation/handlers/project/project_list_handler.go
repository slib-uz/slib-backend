package project

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/projectusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type ProjectsListHandler struct {
	uc *ProjectListUseCase
}

// @inject
func NewProjectsListHandler(uc *ProjectListUseCase) *ProjectsListHandler {
	return &ProjectsListHandler{uc: uc}
}

// Handle
// @Tags project
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {array} entity.ProjectEntity
// @Router /project/list [get]
func (this *ProjectsListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(ctx)
	paging, err := this.uc.Execute(page, pageSize)

	if err != nil {
		return err
	}

	return c.JsonResponse(200, paging)
}
