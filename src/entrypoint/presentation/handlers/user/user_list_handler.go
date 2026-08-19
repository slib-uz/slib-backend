package user

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/userusecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type UserListHandler struct {
	uc *userusecases.UserListUseCase
}

// @inject
func NewUserListHandler(uc *userusecases.UserListUseCase) *UserListHandler {
	return &UserListHandler{uc: uc}
}

// Handle
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page         query     int     false  "Page number"
// @Param        page_size    query     int     false  "Page size"
// @Param        search       query     string  false  "Search query"
// @Success      200  {array} entity.UserEntity
// @Failure      403  {object} response.Response
// @Router       /user/list [get]
func (this *UserListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)
	user := c.User

	page, pageSize := context2.GetPagingParams(c)
	search := c.QueryParam("search")

	result, err := this.uc.Execute(page, pageSize, search, user)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, result)
}
