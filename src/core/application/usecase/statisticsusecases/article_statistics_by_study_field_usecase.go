package statisticsusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ArticleStatisticsByStudyFieldUseCase struct {
	statisticsRepository repository.StatisticsRepository
}

// @inject
func NewArticleStatisticsByStudyFieldUseCase(statisticsRepository repository.StatisticsRepository) *ArticleStatisticsByStudyFieldUseCase {
	return &ArticleStatisticsByStudyFieldUseCase{
		statisticsRepository: statisticsRepository,
	}
}

func (this *ArticleStatisticsByStudyFieldUseCase) Execute(page, pageSize int, journalID *uint) (*entity2.PagingEntity[entity2.ArticleStatisticsByStudyFieldEntity], error) {
	stats, err := this.statisticsRepository.GetArticleStatisticsByStudyField(page, pageSize, journalID)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
