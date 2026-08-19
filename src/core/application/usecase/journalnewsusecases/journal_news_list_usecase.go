package journalnewsusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalNewsListUseCase struct {
	repository repository.JournalNewsRepository
}

// @inject
func NewJournalNewsListUseCase(repository repository.JournalNewsRepository) *JournalNewsListUseCase {
	return &JournalNewsListUseCase{repository: repository}
}

func (this *JournalNewsListUseCase) Execute(journalID uint, page, pageSize int) (*entity.PagingEntity[entity.JournalNewsEntity], error) {
	return this.repository.GetByJournalID(journalID, page, pageSize)
}
