package doiusecases

import (
	"context"

	"slib.uz/src/core/domain/ports/gateway"
)

type CrossRefCheckAuthUseCase struct {
	crossRefGateway gateway.CrossRefGateway
}

// @inject
func NewCrossRefCheckAuthUseCase(crossRefGateway gateway.CrossRefGateway) *CrossRefCheckAuthUseCase {
	return &CrossRefCheckAuthUseCase{crossRefGateway: crossRefGateway}
}

func (this *CrossRefCheckAuthUseCase) Execute(ctx context.Context, username, password string) (bool, error) {
	return this.crossRefGateway.CheckAuth(username, password)
}
