package article_application

import (
	"net/http"

	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/article_applications_usecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type ApplicationDetailHandler struct {
	uc *ApplicationDetailUsecase
}

// @inject
func NewApplicationDetailHandler(uc *ApplicationDetailUsecase) *ApplicationDetailHandler {
	return &ApplicationDetailHandler{uc: uc}
}

// Handle
// @Tags article-application
// @Accept json
// @Produce json
// @Param applicationId path int true "Application ID"
// @Success 200 {object} entity.JournalApplicationEntity
// @Security BearerAuth
// @Router /article-application/detail/{applicationId} [get]
func (this *ApplicationDetailHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	applicationId, err := context.GetIntPathParam(ctx, "applicationId")
	if err != nil {
		return err
	}

	app, err := this.uc.Execute(c.User, uint(applicationId))
	if err != nil {
		return err
	}
	return c.JsonResponse(http.StatusOK, app)
}
