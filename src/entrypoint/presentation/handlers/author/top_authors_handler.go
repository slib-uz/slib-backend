package author

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/authorusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type TopAuthorsHandler struct {
	uc *authorusecases.TopAuthorsListUsecase
}

// @inject
func NewTopAuthorsHandler(uc *authorusecases.TopAuthorsListUsecase) *TopAuthorsHandler {
	return &TopAuthorsHandler{uc: uc}
}

// Handle 	godoc
// @Tags 	author
// @Accept 	json
// @Produce json
// @Param 	top query int false "Count"
// @Param 	journal_id query int false "Filter by journal ID (from current journal domain/tenant)"
// @Success 200 {array} entity.AuthorEntity
// @Router 	/author/top [get]
func (this *TopAuthorsHandler) Handle(ctx echo.Context) error {

	c := ctx.(*context2.Context)

	top := context2.GetIntQueryParam(ctx, "top", 10)
	if top >= 50 {
		return response.NewFailResponse(400, "Top count must be less than 50")
	}

	var journalID *uint
	if val := c.QueryParam("journal_id"); val != "" {
		id, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return response.NewFailResponse(400, "Invalid journal_id")
		}
		rid := uint(id)
		journalID = &rid
	}

	result, err := this.uc.Execute(top, journalID)
	if err != nil {
		return err
	}

	return c.JsonResponse(http.StatusOK, result)
}
