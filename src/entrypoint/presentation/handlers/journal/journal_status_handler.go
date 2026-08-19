package journal

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/journalusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type JournalStatusHandler struct {
	uc *journalusecases.ChangeJournalStatusUseCase
}

// @inject
func NewJournalStatusHandler(uc *journalusecases.ChangeJournalStatusUseCase) *JournalStatusHandler {
	return &JournalStatusHandler{uc: uc}
}

type JournalStatusRequest struct {
	IsActive bool `json:"is_active"`
}

// Handle JournalStatusHandler
// @Tags         journal-manage
// @Summary      Change journal IsActive status
// @Description  Allows admins to enable or disable a journal
// @Accept       json
// @Produce      json
// @Param        id      path     int                   true "Journal ID"
// @Param        body    body     JournalStatusRequest  true "Status body"
// @Success      200     {object} map[string]string
// @Security     BearerAuth
// @Router       /journal-manage/{id}/status [patch]
func (this *JournalStatusHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)
	idStr := c.Param("id")
	journalID64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return err
	}
	journalID := uint(journalID64)

	var req JournalStatusRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := this.uc.Execute(c.User, journalID, req.IsActive); err != nil {
		return err
	}

	return c.JsonResponse(200, map[string]string{"message": "Journal status updated"})
}
