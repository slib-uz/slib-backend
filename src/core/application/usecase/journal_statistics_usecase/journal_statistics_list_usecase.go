package journalstatisticsusecase

import (
	"slib.uz/src/core/application/usecase/journalusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalStatisticListUseCase struct {
	journalRepository        repository.JournalRepository
	completionPercentUseCase *journalusecases.JournalCompletionPercentUseCase
}

// @inject
func NewJournalStatisticListUseCase(
	journalRepository repository.JournalRepository,
	completionPercentUseCase *journalusecases.JournalCompletionPercentUseCase,
) *JournalStatisticListUseCase {
	return &JournalStatisticListUseCase{
		journalRepository:        journalRepository,
		completionPercentUseCase: completionPercentUseCase,
	}
}

func (this *JournalStatisticListUseCase) Execute(page, pageSize int, institutionID, publisherID uint, name, description, issn, publisherName *string) (*entity.PagingEntity[entity.JournalStatisticV2Entity], error) {
	statistics, err := this.journalRepository.GetJournalStatisticsV2(page, pageSize, institutionID, publisherID, name, description, issn, publisherName)
	if err != nil {
		return nil, err
	}

	journalIDs := make([]uint, len(statistics.Items))
	for i, stat := range statistics.Items {
		journalIDs[i] = stat.JournalID
	}

	completionPercentMap, err := this.completionPercentUseCase.ExecuteBatch(journalIDs)
	if err != nil {
		return nil, err
	}

	for i, stat := range statistics.Items {
		if percent, exists := completionPercentMap[stat.JournalID]; exists {
			statistics.Items[i].CompletionPercent = percent
		}
	}

	return statistics, nil
}
