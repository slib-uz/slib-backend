package journalratingusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalRatingStatsUseCase struct {
	repository repository.JournalRatingRepository
}

// @inject
func NewJournalRatingStatsUseCase(repository repository.JournalRatingRepository) *JournalRatingStatsUseCase {
	return &JournalRatingStatsUseCase{repository: repository}
}

func (this *JournalRatingStatsUseCase) Execute(journalID uint) (*entity.JournalRatingStatsEntity, error) {
	stats, err := this.repository.GetStatsByJournalID(journalID)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
