package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/chief_editor"
)

type ChiefEditorGroup struct {
	chiefEditorByJournalIDHandler *chief_editor.ChiefEditorByJournalIDHandler
}

// @inject
func NewChiefEditorGroup(chiefEditorByJournalIDHandler *chief_editor.ChiefEditorByJournalIDHandler) *ChiefEditorGroup {
	return &ChiefEditorGroup{chiefEditorByJournalIDHandler: chiefEditorByJournalIDHandler}
}

func (this *ChiefEditorGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/:journal_id", this.chiefEditorByJournalIDHandler.Handle)
}
