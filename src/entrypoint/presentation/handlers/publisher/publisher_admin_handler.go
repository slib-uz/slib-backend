package publisher

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	. "slib.uz/src/core/application/usecase/publisherusecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type PublisherAdminHandler struct {
	uc *PublisherAdminUseCase
}

// @inject
func NewPublisherAdminHandler(uc *PublisherAdminUseCase) *PublisherAdminHandler {
	return &PublisherAdminHandler{uc: uc}
}

// Handle
// @Tags Publisher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param publisherId path int true "Publisher ID"
// @Success 200 {array} entity.UserRoleWithBasicUserEntity
// @Failure 400 {object} response.Response
// @Router /publisher/{publisherId}/admins [get]
func (this *PublisherAdminHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)
	result, err := this.uc.Execute(this.GetPublisherID(ctx))
	if err != nil {
		return err
	}
	return c.JsonResponse(200, result)
}

func (this *PublisherAdminHandler) GetPublisherID(ctx echo.Context) uint {
	publisherId := ctx.Param("publisherId")
	if publisherId == "" {
		panic(response.NewFailResponse(400, "Publisher ID is required"))
	}
	id, err := strconv.Atoi(publisherId)
	if err != nil {
		panic(response.NewFailResponse(400, "Invalid Publisher ID"))
	}
	return uint(id)
}
