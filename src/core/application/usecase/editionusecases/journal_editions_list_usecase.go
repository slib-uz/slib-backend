package editionusecases

import (
	"context"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalEditionsListUseCase struct {
	repository repository.EditionRepository
}

// @inject
func NewJournalEditionsListUseCase(repository repository.EditionRepository) *JournalEditionsListUseCase {
	return &JournalEditionsListUseCase{repository: repository}
}

func (this *JournalEditionsListUseCase) Execute(ctx context.Context, journalID uint, page, pageSize int, search string, year int) (*entity.PagingEntity[entity.EditionEntity], error) {
	return this.repository.GetByJournalID(ctx, journalID, page, pageSize, search, year)
}
