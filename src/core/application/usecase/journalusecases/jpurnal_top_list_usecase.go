package journalusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type TopJournalListUseCase struct {
	repository repository.JournalRepository
}

// @inject
func NewTopJournalListUseCase(repository repository.JournalRepository) *TopJournalListUseCase {
	return &TopJournalListUseCase{repository: repository}
}

func (this *TopJournalListUseCase) Execute(page, pageSize int) (*entity2.PagingEntity[entity2.JournalEntity], error) {
	paging, err := this.repository.GetTopJournals(page, pageSize)
	if err != nil {
		return nil, err
	}
	return entity2.NewPagingEntity(page, pageSize, paging.Total, paging.Items), nil
}
