package author

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	. "slib.uz/src/core/application/usecase/authorusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type AuthorsListHandler struct {
	uc *AuthorsListUsecase
}

// @inject
func NewAuthorsListHandler(uc *AuthorsListUsecase) *AuthorsListHandler {
	return &AuthorsListHandler{uc: uc}
}

// Handle
// @Tags author
// @Accept  json
// @Produce  json
// @Param name query string false "Author name"
// @Param science_id query string false "Science ID"
// @Param journal_id query int false "Filter by journal ID (from current journal domain/tenant)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {array} entity.AuthorEntity
// @Router /author/list [get]
func (this *AuthorsListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(ctx)
	name := c.QueryParam("name")
	scienceID := c.QueryParam("science_id")

	var journalID *uint
	if val := c.QueryParam("journal_id"); val != "" {
		id, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return response.NewFailResponse(400, "Invalid journal_id")
		}
		rid := uint(id)
		journalID = &rid
	}

	paging, err := this.uc.Execute(page, pageSize, name, scienceID, journalID)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, paging)
}
