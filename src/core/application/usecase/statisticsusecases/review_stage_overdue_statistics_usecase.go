package statisticsusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ReviewStageOverdueStatisticsUseCase struct {
	reviewStageRepository repository.ReviewStageRepository
}

// @inject
func NewReviewStageOverdueStatisticsUseCase(reviewStageRepository repository.ReviewStageRepository) *ReviewStageOverdueStatisticsUseCase {
	return &ReviewStageOverdueStatisticsUseCase{
		reviewStageRepository: reviewStageRepository,
	}
}

func (this *ReviewStageOverdueStatisticsUseCase) Execute(journalID *uint) ([]*entity.ReviewStageOverdueStatisticsEntity, error) {
	return this.reviewStageRepository.GetOverdueStatistics(journalID)
}
