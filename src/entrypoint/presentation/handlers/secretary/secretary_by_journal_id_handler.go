package secretary

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/secretaryusecase"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type SecretaryByJournalIDHandler struct {
	uc *SecretaryByJournalIDUseCase
}

// @inject
func NewSecretaryByJournalIDHandler(uc *SecretaryByJournalIDUseCase) *SecretaryByJournalIDHandler {
	return &SecretaryByJournalIDHandler{uc: uc}
}

// Handle SecretaryByJournalIDHandler
// @Tags secretary
// @Accept json
// @Produce json
// @Param journal_id path int true "Journal ID"
// @Success 200 {array} entity.SecretaryEntity
// @Failure 400 {object} response.Response
// @Router /secretary/{journal_id} [get]
func (this *SecretaryByJournalIDHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	journalId, err := context2.GetIntPathParam(c, "journal_id")
	if err != nil {
		return err
	}

	secretaries, err := this.uc.Execute(uint(journalId))
	if err != nil {
		return err
	}
	return c.JsonResponse(200, secretaries)
}
