package journalstatisticsusecase

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/journalusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalStatisticsDetailUseCase struct {
	journalRepository        repository.JournalRepository
	completionPercentUseCase *journalusecases.JournalCompletionPercentUseCase
}

// @inject
func NewJournalStatisticsDetailUseCase(
	journalRepository repository.JournalRepository,
	completionPercentUseCase *journalusecases.JournalCompletionPercentUseCase,
) *JournalStatisticsDetailUseCase {
	return &JournalStatisticsDetailUseCase{
		journalRepository:        journalRepository,
		completionPercentUseCase: completionPercentUseCase,
	}
}

func (this *JournalStatisticsDetailUseCase) Execute(journalID uint) (*entity.JournalStatisticV2Entity, error) {
	statistic, err := this.journalRepository.GetJournalStatisticV2ByJournalID(journalID)
	if err != nil {
		return nil, err
	}
	if statistic == nil {
		return nil, response.NotFoundError
	}

	completionPercentMap, err := this.completionPercentUseCase.ExecuteBatch([]uint{journalID})
	if err != nil {
		return nil, err
	}

	if percent, exists := completionPercentMap[journalID]; exists {
		statistic.CompletionPercent = percent
	}

	return statistic, nil
}
