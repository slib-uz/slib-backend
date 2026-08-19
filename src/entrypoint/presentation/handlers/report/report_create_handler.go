package report

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/reportuscases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/handlers/report/schema"
)

type ReportCreateHandler struct {
	uc *ReportCreateUseCase
}

// @inject
func NewReportCreateHandler(uc *ReportCreateUseCase) *ReportCreateHandler {
	return &ReportCreateHandler{uc: uc}
}

// Handle ReportCreateHandler
// @Tags report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param report body schema.ReportCreateRequest true "ReportCreateRequest"
// @Success 201 {object} response.Response
// @Router /report/create [post]
func (this ReportCreateHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	data, err := context2.GetBody[schema.ReportCreateRequest](c)
	if err != nil {
		return err
	}

	if err := this.uc.Execute(c.User.ID, data.ToEntity()); err != nil {
		return c.JsonResponse(400, err.Error())
	}

	return c.JsonResponse(201, "Your report has been sent to admin")
}
