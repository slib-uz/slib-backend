package journal

import (
	"net/http"

	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/journalusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type JournalReviewersListHandler struct {
	uc *JournalReviewersListUsecase
}

// @inject
func NewJournalReviewersListHandler(uc *JournalReviewersListUsecase) *JournalReviewersListHandler {
	return &JournalReviewersListHandler{uc: uc}
}

// Handle
// @Tags         journal
// @Param        journalId path int true "Journal ID"
// @Accept       json
// @Produce      json
// @Success      200 {array} entity.ReviewerEntity
// @Router       /journal/{journalId}/reviewers [get]
func (this *JournalReviewersListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	journalId, err := context2.GetIntPathParam(ctx, "journalId")
	if err != nil {
		return err
	}

	reviewers, err := this.uc.Execute(uint(journalId))
	if err != nil {
		return err
	}

	return c.JsonResponse(http.StatusOK, reviewers)
}
