package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/aidetect"
)

type AiDetectGroup struct {
	aiDetectResultsHandler *aidetect.AiDetectResultsHandler
}

// @inject
func NewAiDetectGroup(aiDetectResultsHandler *aidetect.AiDetectResultsHandler) *AiDetectGroup {
	return &AiDetectGroup{aiDetectResultsHandler: aiDetectResultsHandler}
}

func (this *AiDetectGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/results", this.aiDetectResultsHandler.Handle)
}
