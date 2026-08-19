package publisher

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/publisherusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type PublisherListHandler struct {
	uc *publisherusecases.PublisherListUseCase
}

// @inject
func NewPublisherListHandler(uc *publisherusecases.PublisherListUseCase) *PublisherListHandler {
	return &PublisherListHandler{uc: uc}
}

// Handle handles the request to get a list of publishers
// @Tags Publisher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param name query string false "Query name"
// @Param tin query string false "Query TIN"
// @Param institution_id query int false "Institution ID"
// @Param unassigned query bool false "Filter publishers not assigned to the institution"
// @Success 200 {array} entity.PublisherEntity
// @Router /publisher/list [get]
func (this *PublisherListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, size := context2.GetPagingParams(ctx)
	name := ctx.QueryParam("name")
	tin := ctx.QueryParam("tin")
	institutionID := uint(context2.GetIntQueryParam(ctx, "institution_id", 0))
	unassigned := ctx.QueryParam("unassigned") == "true"

	paging, err := this.uc.Execute(page, size, name, tin, institutionID, unassigned)

	if err != nil {
		return err
	}

	return c.JsonResponse(200, paging)

}

// func (this *PublisherListHandler) getSearchParams(c echo.Context) (name *string, tin *string) {
// 	qname, qtin := c.QueryParam("name"), c.QueryParam("tin")
// 	if qname != "" {
// 		name = &qname
// 	}
// 	if qtin != "" {
// 		tin = &qtin
// 	}
// 	return name, tin
// }
