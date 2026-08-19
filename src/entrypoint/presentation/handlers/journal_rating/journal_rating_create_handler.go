package journal_rating

import (
	"github.com/labstack/echo/v4"
	ratingusecases "slib.uz/src/core/application/usecase/journalratingusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/journal_rating/schema"
)

type JournalRatingCreateHandler struct {
	uc *ratingusecases.JournalRatingCreateUseCase
}

// @inject
func NewJournalRatingCreateHandler(uc *ratingusecases.JournalRatingCreateUseCase) *JournalRatingCreateHandler {
	return &JournalRatingCreateHandler{uc: uc}
}

// Handle JournalRatingCreateHandler
// @Tags journal-rating
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param journal-rating body schema.JournalRatingCreateRequest true "JournalRatingCreateRequest"
// @Success 200 {object} response.Response
// @Router /journal-rating/create [post]
func (this *JournalRatingCreateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	data, err := context2.GetBody[schema.JournalRatingCreateRequest](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User.ID, data.ToEntity()); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	return c.JsonResponse(200, "Journal rating created successfully")
}
