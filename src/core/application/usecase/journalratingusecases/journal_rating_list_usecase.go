package journalratingusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalRatingListUseCase struct {
	repository repository.JournalRatingRepository
}

// @inject
func NewJournalRatingListUseCase(repository repository.JournalRatingRepository) *JournalRatingListUseCase {
	return &JournalRatingListUseCase{repository: repository}
}

func (this *JournalRatingListUseCase) Execute(journalID uint, page, pageSize int, ordering string) (*entity2.PagingEntity[entity2.JournalRatingEntity], error) {
	paging, err := this.repository.GetByJournalID(journalID, page, pageSize, ordering)

	if err != nil {
		return nil, err
	}

	return paging, nil
}
