package project

import (
	"strconv"

	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/projectusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type ProjectDetailHandler struct {
	uc *ProjectDetailUseCase
}

// @inject
func NewProjectDetailHandler(uc *ProjectDetailUseCase) *ProjectDetailHandler {
	return &ProjectDetailHandler{uc: uc}
}

// Handle
// @Tags project
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} entity.ProjectEntity
// @Failure 404 {object} response.Response
// @Router /project/detail/{id} [get]
func (this *ProjectDetailHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JsonResponse(400, "Invalid ID")
	}

	project, err := this.uc.Execute(uint(id))
	if err != nil {
		return err
	}

	return c.JsonResponse(200, project)
}
