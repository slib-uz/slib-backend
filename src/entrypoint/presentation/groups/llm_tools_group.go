package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/llm_tools"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
)

type LlmToolsGroup struct {
	extractArticleMetadataHandler *llm_tools.ExtractArticleMetadataHandler
}

// @inject
func NewLlmToolsGroup(extractArticleMetadataHandler *llm_tools.ExtractArticleMetadataHandler) *LlmToolsGroup {
	return &LlmToolsGroup{extractArticleMetadataHandler: extractArticleMetadataHandler}
}

func (this *LlmToolsGroup) RegisterRoutes(group *echo.Group) {
	authenticated := group.Group("", permissions.AuthenticatedPermission)
	authenticated.POST("/article-metadata", this.extractArticleMetadataHandler.Handle)
}
