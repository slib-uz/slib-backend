package journal_rating

import (
	"github.com/labstack/echo/v4"
	ratingusecases "slib.uz/src/core/application/usecase/journalratingusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type JournalRatingDeleteHandler struct {
	uc *ratingusecases.JournalRatingDeleteUseCase
}

// @inject
func NewJournalRatingDeleteHandler(uc *ratingusecases.JournalRatingDeleteUseCase) *JournalRatingDeleteHandler {
	return &JournalRatingDeleteHandler{uc: uc}
}

// Handle JournalRatingDeleteHandler
// @Tags journal-rating
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Journal Rating ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /journal-rating/{id} [delete]
func (this *JournalRatingDeleteHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	id, err := context2.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}
	if err := this.uc.Execute(c.User.ID, uint(id)); err != nil {
		return c.JsonResponse(400, err.Error())
	}
	return c.JsonResponse(200, "Journal rating deleted successfully")
}
