package uzsciusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
)

type UzSciFormListUseCase struct {
	uzsciGateway gateway.UzsciGateway
}

// @inject
func NewUzSciFormListUseCase(uzsciGateway gateway.UzsciGateway) *UzSciFormListUseCase {
	return &UzSciFormListUseCase{uzsciGateway: uzsciGateway}
}

func (this *UzSciFormListUseCase) Execute(ratingPeriodID uint) ([]*entity.UzSciFormEntity, error) {
	return this.uzsciGateway.GetAllForms(ratingPeriodID)
}
