package antiplagusecases

import (
	"slib.uz/src/core/domain/ports/gateway"
)

type BalanceUseCase struct {
	gateway gateway.AntiPlagGateway
}

// @inject
func NewBalanceUseCase(gateway gateway.AntiPlagGateway) *BalanceUseCase {
	return &BalanceUseCase{gateway: gateway}
}

func (uc *BalanceUseCase) Execute() (int, error) {
	return uc.gateway.GetBalance()
}
