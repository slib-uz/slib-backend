package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/ethics_policy"
)

type EthicsPolicyGroup struct {
	EthicsPolicyListHandler   *ethics_policy.EthicsPolicyListHandler
	EthicsPolicyDetailHandler *ethics_policy.EthicsPolicyDetailHandler
}

// @inject
func NewEthicsPolicyGroup(
	ethicsPolicyListHandler *ethics_policy.EthicsPolicyListHandler,
	ethicsPolicyDetailHandler *ethics_policy.EthicsPolicyDetailHandler) *EthicsPolicyGroup {
	return &EthicsPolicyGroup{
		EthicsPolicyListHandler:   ethicsPolicyListHandler,
		EthicsPolicyDetailHandler: ethicsPolicyDetailHandler,
	}
}

func (this *EthicsPolicyGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/list", this.EthicsPolicyListHandler.Handle)
	group.GET("/detail/:id", this.EthicsPolicyDetailHandler.Handle)
}
