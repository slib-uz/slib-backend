package journal_config

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/journalconfigusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type JournalConfigHandler struct {
	us *JournalConfigUseCase
}

// @inject
func NewJournalConfigHandler(us *JournalConfigUseCase) *JournalConfigHandler {
	return &JournalConfigHandler{us: us}
}

// Handle Get Journal Config
// @Tags journal-config
// @Accept: application/json
// @Produce: application/json
// @Param journal_id query int false "Journal ID"
// @Param website_url query string false "Website URL"
// @Success 200 {object} entity.JournalConfigEntity
// @Failure 404 {object} response.Response
// @Router /journal-config [get]
func (this *JournalConfigHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	journalID := context.GetIntQueryParam(c, "journal_id", 0)

	websiteURL := c.QueryParam("website_url")

	entity, err := this.us.Execute(uint(journalID), websiteURL, c.User)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, entity)
}
