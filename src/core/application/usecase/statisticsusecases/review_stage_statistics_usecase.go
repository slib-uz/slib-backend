package statisticsusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type ReviewStageStatisticsUseCase struct {
	reviewStageRepository repository.ReviewStageRepository
}

// @inject
func NewReviewStageStatisticsUseCase(reviewStageRepository repository.ReviewStageRepository) *ReviewStageStatisticsUseCase {
	return &ReviewStageStatisticsUseCase{
		reviewStageRepository: reviewStageRepository,
	}
}

func (this *ReviewStageStatisticsUseCase) Execute(journalID *uint) ([]*entity.ReviewStageStatisticsEntity, error) {
	stats, err := this.reviewStageRepository.GetStatistics(journalID)
	if err != nil {
		return nil, err
	}

	overdueStats, err := this.reviewStageRepository.GetOverdueStatistics(journalID)
	if err != nil {
		return nil, err
	}

	overdueByStage := make(map[enum.Stage]int64, len(overdueStats))
	for _, overdue := range overdueStats {
		overdueByStage[overdue.StageNumber] = overdue.OverdueCount
	}

	for _, stat := range stats {
		stat.OverdueCount = overdueByStage[stat.StageNumber]
	}

	if stats == nil {
		stats = append(stats, entity.NewReviewStageStatisticsEntity("Technical Review", 10, 0, 0, 0, 0, 0))
		stats = append(stats, entity.NewReviewStageStatisticsEntity("Peer Review", 20, 0, 0, 0, 0, 0))
		stats = append(stats, entity.NewReviewStageStatisticsEntity("Payment", 30, 0, 0, 0, 0, 0))
		stats = append(stats, entity.NewReviewStageStatisticsEntity("Publish", 40, 0, 0, 0, 0, 0))
	}

	return stats, nil
}
