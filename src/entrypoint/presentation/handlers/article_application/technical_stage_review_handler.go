package article_application

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/article_applications_usecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/article_application/schema"
)

type TechnicalStageReviewHandler struct {
	uc *TechnicalReviewUsecase
}

// @inject
func NewTechnicalStageReviewHandler(uc *TechnicalReviewUsecase) *TechnicalStageReviewHandler {
	return &TechnicalStageReviewHandler{uc: uc}
}

// Handle
// @Tags article-application
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body schema.TechnicalReviewRequest true "Technical review request"
// @Success 200 {object} response.Response
// @Router /article-application/review/technical-review-stage [post]
func (this *TechnicalStageReviewHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	data, err := context2.GetBody[schema.TechnicalReviewRequest](ctx)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User, data.ApplicationID, data.Status, data.Reason); err != nil {
		return err
	}

	return c.JsonResponse(200, "Technical review status updated successfully")

}
