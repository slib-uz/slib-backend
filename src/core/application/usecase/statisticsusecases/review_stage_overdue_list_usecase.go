package statisticsusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ReviewStageOverdueListUseCase struct {
	reviewStageRepository repository.ReviewStageRepository
}

// @inject
func NewReviewStageOverdueListUseCase(reviewStageRepository repository.ReviewStageRepository) *ReviewStageOverdueListUseCase {
	return &ReviewStageOverdueListUseCase{
		reviewStageRepository: reviewStageRepository,
	}
}

func (this *ReviewStageOverdueListUseCase) Execute(page, pageSize int, journalID *uint) (*entity.PagingEntity[entity.ReviewStageOverdueEntity], error) {
	return this.reviewStageRepository.GetOverdueList(page, pageSize, journalID)
}
