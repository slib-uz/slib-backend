package soato

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/soatousecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type SoatoListHandler struct {
	uc *soatousecases.SoatoListUseCase
}

// @inject
func NewSoatoListHandler(uc *soatousecases.SoatoListUseCase) *SoatoListHandler {
	return &SoatoListHandler{uc: uc}
}

// Handle
// @Tags Soato
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param is_children query bool false "Is Children (true: children with parents, false: parents with children)"
// @Success 200 {array} entity.SoatoClassifierEntity
// @Router /soato/list [get]
func (this *SoatoListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, size := context2.GetPagingParams(c)

	isChildrenParam := c.QueryParam("is_children")
	isChildren := false
	if isChildrenParam != "" {
		var err error
		isChildren, err = strconv.ParseBool(isChildrenParam)
		if err != nil {
			// default to false if invalid
			isChildren = false
		}
	}

	result, err := this.uc.Execute(page, size, isChildren)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
