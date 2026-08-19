package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/secretary"
)

type SecretaryGroup struct {
	secretaryByJournalIDHandler *secretary.SecretaryByJournalIDHandler
}

// @inject
func NewSecretaryGroup(secretaryByJournalIDHandler *secretary.SecretaryByJournalIDHandler) *SecretaryGroup {
	return &SecretaryGroup{secretaryByJournalIDHandler: secretaryByJournalIDHandler}
}

func (this *SecretaryGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/:journal_id", this.secretaryByJournalIDHandler.Handle)
}
