package public_offer

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/public_offer_usecases"
	"slib.uz/src/entrypoint/presentation/app/context"
)

type PublicOfferHandler struct {
	uc *PublicOfferUseCase
}

// @inject
func NewPublicOfferHandler(uc *PublicOfferUseCase) *PublicOfferHandler {
	return &PublicOfferHandler{uc: uc}
}

// Handle
// @Summary      Get public offer
// @Description  Get the public offer document
// @Tags         public-offer
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.PublicOfferEntity
// @Failure      400  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /public-offer [get]
func (this *PublicOfferHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context.Context)

	publicOffer, err := this.uc.Execute()
	if err != nil {
		return err
	}

	return c.JsonResponse(200, publicOffer)
}
