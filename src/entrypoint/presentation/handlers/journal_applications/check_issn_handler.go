package journal_applications

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/journal_applications_usecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type CheckISSNHandler struct {
	uc *CheckISSNUsecase
}

// @inject
func NewCheckISSNHandler(uc *CheckISSNUsecase) *CheckISSNHandler {
	return &CheckISSNHandler{uc: uc}
}

// Handle
// @Tags journal-applications
// @Accept: application/json
// @Produce: application/json
// @Param issn query string true "ISSNPaper"
// @Param issn_type query string true "ISSN Type (issn_paper or issn_online)"
// @Success 200 {object} entity.JournalApplicationEntity
// @Failure 404 {object} response.Response
// @Router /journal-applications/check-issn [get]
func (this *CheckISSNHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	issn := c.QueryParam("issn")
	issnType := c.QueryParam("issn_type")

	if issnType == "" || (issnType != "issn_paper" && issnType != "issn_online") {
		return c.JsonResponse(400, "Invalid ISSN type. Must be 'issn_paper' or 'issn_online'.")
	}

	if err := this.uc.Execute(issn, issnType); err != nil {
		return err
	}
	return c.JsonResponse(404, "Not Found")

}
