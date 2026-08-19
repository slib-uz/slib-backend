package editionusecases

import (
	"context"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type EditionDetailUseCase struct {
	repository repository.EditionRepository
}

// @inject
func NewEditionDetailUseCase(repository repository.EditionRepository) *EditionDetailUseCase {
	return &EditionDetailUseCase{repository: repository}
}

func (this *EditionDetailUseCase) Execute(ctx context.Context, id uint) (*entity.EditionEntity, error) {
	return this.repository.GetByID(ctx, id)
}
