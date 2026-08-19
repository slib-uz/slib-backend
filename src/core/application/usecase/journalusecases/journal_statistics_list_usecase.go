package journalusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalStatisticsListUseCase struct {
	journalRepository repository.JournalRepository
}

// @inject
func NewJournalStatisticsListUseCase(journalRepository repository.JournalRepository) *JournalStatisticsListUseCase {
	return &JournalStatisticsListUseCase{
		journalRepository: journalRepository,
	}
}

func (this *JournalStatisticsListUseCase) Execute(page, pageSize int) (*entity.PagingEntity[entity.JournalStatisticEntity], error) {
	return this.journalRepository.GetJournalStatistics(page, pageSize)
}
