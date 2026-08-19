package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/instruction"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
)

type InstructionGroup struct {
	createOrUpdateHandler *instruction.InstructionCreateOrUpdateHandler
	deleteHandler         *instruction.InstructionDeleteHandler
	listHandler           *instruction.InstructionListHandler
	detailHandler         *instruction.InstructionDetailHandler
}

// @inject
func NewInstructionGroup(
	createOrUpdateHandler *instruction.InstructionCreateOrUpdateHandler,
	deleteHandler *instruction.InstructionDeleteHandler,
	listHandler *instruction.InstructionListHandler,
	detailHandler *instruction.InstructionDetailHandler,
) *InstructionGroup {
	return &InstructionGroup{
		createOrUpdateHandler: createOrUpdateHandler,
		deleteHandler:         deleteHandler,
		listHandler:           listHandler,
		detailHandler:         detailHandler,
	}
}

func (this *InstructionGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/save", this.createOrUpdateHandler.Handle, permissions.AdminPermission)
	group.DELETE("/:id/delete", this.deleteHandler.Handle, permissions.AdminPermission)
	group.GET("/list", this.listHandler.Handle)
	group.GET("/:key", this.detailHandler.Handle)
}
