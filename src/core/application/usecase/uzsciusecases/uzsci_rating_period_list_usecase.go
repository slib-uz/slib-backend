package uzsciusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
)

type UzSciRatingPeriodListUseCase struct {
	uzsciGateway gateway.UzsciGateway
}

// @inject
func NewUzSciRatingPeriodListUseCase(uzsciGateway gateway.UzsciGateway) *UzSciRatingPeriodListUseCase {
	return &UzSciRatingPeriodListUseCase{uzsciGateway: uzsciGateway}
}

func (this *UzSciRatingPeriodListUseCase) Execute(isActive *bool) ([]*entity.UzSciRatingPeriodEntity, error) {
	return this.uzsciGateway.GetPublicRatingPeriods(isActive)
}
