package journal_applications

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/journal_applications_usecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/journal_applications/schema"
)

type ReviewApplicationHandler struct {
	uc *ReviewApplicationUseCase
}

// @inject
func NewReviewApplicationHandler(uc *ReviewApplicationUseCase) *ReviewApplicationHandler {
	return &ReviewApplicationHandler{uc: uc}
}

// Handle ReviewApplicationHandler
// @Tags journal-applications
// @Accept json
// @Produce json
// @Param request body schema.ReviewApplicationRequest true "UpdateToNext Application"
// @Success 200 {object} schema.ReviewApplicationRequest
// @Security BearerAuth
// @Router /journal-applications/review [post]
func (this *ReviewApplicationHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)
	data, err := context2.GetBody[schema.ReviewApplicationRequest](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(data.ApplicationID, data.Status, data.RejectionReason); err != nil {
		return err
	}
	return c.JsonResponse(200, data)

}
