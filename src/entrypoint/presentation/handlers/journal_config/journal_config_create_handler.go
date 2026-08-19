package journal_config

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/journalconfigusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/journal_config/schema"
)

type JournalConfigCreateHandler struct {
	uc *CreateOrUpdateJournalConfigUseCase
}

// @inject
func NewJournalConfigCreateHandler(uc *CreateOrUpdateJournalConfigUseCase) *JournalConfigCreateHandler {
	return &JournalConfigCreateHandler{uc: uc}
}

// Handle Create Or Update Journal Config
// @Tags journal-config
// @Security BearerAuth
// @Accept: application/json
// @Produce: application/json
// @Param body body schema.JournalConfigSchema true "Journal Config"
// @Success 200 {object} entity.JournalConfigEntity
// @Failure 404 {object} response.Response
// @Router /journal-config/create-or-update [post]
func (this *JournalConfigCreateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)
	body, err := context.GetBody[schema.JournalConfigSchema](c)
	if err != nil {
		return err
	}

	err = this.uc.Execute(body.ToEntity(c.User.ID), c.User)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, body)
}
