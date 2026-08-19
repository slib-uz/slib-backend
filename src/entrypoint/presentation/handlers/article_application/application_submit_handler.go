package article_application

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	. "slib.uz/src/core/application/usecase/article_applications_usecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/article_application/schema"
)

type ApplicationSubmitHandler struct {
	uc *ApplicationSubmitUseCase
}

// @inject
func NewApplicationSubmitHandler(uc *ApplicationSubmitUseCase) *ApplicationSubmitHandler {
	return &ApplicationSubmitHandler{uc: uc}
}

// Handle
// @Tags article-application
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body schema.ApplicationSubmitRequest true "Article publish request"
// @Success 201 {object} map[string]any
// @Param journalId path int true "Journal ID"
// @Router /article-application/submit/{journalId} [post]
func (this *ApplicationSubmitHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	data, err := context2.GetBody[schema.ApplicationSubmitRequest](ctx)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(data.ToEntity(), this.GetJournalId(c), c.User.ID); err != nil {
		return err
	}

	return c.JsonResponse(http.StatusCreated, nil)
}

func (this *ApplicationSubmitHandler) GetJournalId(ctx echo.Context) uint {

	journalIdUint, err := strconv.ParseUint(ctx.Param("journalId"), 10, 64)
	if err != nil {
		panic(response.NewFailResponse(400, fmt.Sprint("Invalid journalId: ", err)))
	}
	return uint(journalIdUint)

}
