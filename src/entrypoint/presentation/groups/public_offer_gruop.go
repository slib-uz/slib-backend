package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/public_offer"
)

type PublicOfferGroup struct {
	publicOfferHandler *public_offer.PublicOfferHandler
}

// @inject
func NewPublicOfferGroup(publicOfferHandler *public_offer.PublicOfferHandler) *PublicOfferGroup {
	return &PublicOfferGroup{publicOfferHandler: publicOfferHandler}
}

func (this *PublicOfferGroup) RegisterRoutes(group *echo.Group) {
	group.GET("", this.publicOfferHandler.Handle)
}
