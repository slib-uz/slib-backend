package doiusecases

import (
	"context"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
)

type CrossRefPrefixInfoUseCase struct {
	crossRefGateway gateway.CrossRefGateway
}

// @inject
func NewCrossRefPrefixInfoUseCase(crossRefGateway gateway.CrossRefGateway) *CrossRefPrefixInfoUseCase {
	return &CrossRefPrefixInfoUseCase{crossRefGateway: crossRefGateway}
}

func (this *CrossRefPrefixInfoUseCase) Execute(ctx context.Context, prefix string) (*entity.CrossRefPrefixInfoEntity, error) {
	info, err := this.crossRefGateway.GetPrefixInfo(prefix)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, response.NotFoundError
	}
	return info, nil
}
