package journaleditorialusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalEditorialListUseCase struct {
	repository repository.JournalEditorialRepository
}

// @inject
func NewJournalEditorialListUseCase(repository repository.JournalEditorialRepository) *JournalEditorialListUseCase {
	return &JournalEditorialListUseCase{repository: repository}
}

func (this *JournalEditorialListUseCase) Execute(journalID uint, page, pageSize int) (*entity.PagingEntity[entity.JournalEditorialEntity], error) {
	return this.repository.GetByJournalID(journalID, page, pageSize)
}
