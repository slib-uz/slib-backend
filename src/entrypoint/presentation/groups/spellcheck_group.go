package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/spellcheck"
)

type SpellCheckGroup struct {
	spellCheckResultsListHandler *spellcheck.SpellCheckResultsHandler
}

// @inject
func NewSpellCheckGroup(spellCheckResultsListHandler *spellcheck.SpellCheckResultsHandler) *SpellCheckGroup {
	return &SpellCheckGroup{spellCheckResultsListHandler: spellCheckResultsListHandler}
}

func (this *SpellCheckGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/results", this.spellCheckResultsListHandler.Handle)
}
