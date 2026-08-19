package myapplication

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/article_applications_usecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/myapplication/schema"
)

type ResubmitForTechnicalReviewHandler struct {
	uc *ApplicationResubmitUseCase
}

// @inject
func NewResubmitForTechnicalReviewHandler(uc *ApplicationResubmitUseCase) *ResubmitForTechnicalReviewHandler {
	return &ResubmitForTechnicalReviewHandler{uc: uc}
}

// Handle
// @Security BearerAuth
// @Tags article-application
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Param applicationId path int true "Article ID"
// @Param request body schema.ArticleUpdateRequest true "Article update request"
// @Router /my-application/resubmit-for-technical-review/{applicationId} [post]
func (this *ResubmitForTechnicalReviewHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	data, err := context2.GetBody[schema.ArticleUpdateRequest](ctx)
	if err != nil {
		return err
	}

	applicationId, err := context2.GetIntPathParam(c, "applicationId")
	if err != nil {
		return err
	}

	if err := this.uc.Execute(uint(applicationId), data.ToEntity(), c.User.ID); err != nil {
		return err
	}

	return c.JsonResponse(200, "Application resubmitted successfully")
}
