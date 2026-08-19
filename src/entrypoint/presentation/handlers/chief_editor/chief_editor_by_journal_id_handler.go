package chief_editor

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/chiefeditorusecase"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type ChiefEditorByJournalIDHandler struct {
	uc *ChiefEditorByJournalIDUseCase
}

// @inject
func NewChiefEditorByJournalIDHandler(uc *ChiefEditorByJournalIDUseCase) *ChiefEditorByJournalIDHandler {
	return &ChiefEditorByJournalIDHandler{uc: uc}
}

// Handle ChiefEditorByJournalIDHandler
// @Tags chief editor
// @Accept json
// @Produce json
// @Param journal_id path int true "Journal ID"
// @Success 200 {array} entity.ChiefEditorEntity
// @Failure 400 {object} response.Response
// @Router /chief-editor/{journal_id} [get]
func (this *ChiefEditorByJournalIDHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	journalId, err := context2.GetIntPathParam(c, "journal_id")
	if err != nil {
		return err
	}

	chiefEditors, err := this.uc.Execute(uint(journalId))
	if err != nil {
		return err
	}
	return c.JsonResponse(200, chiefEditors)
}
