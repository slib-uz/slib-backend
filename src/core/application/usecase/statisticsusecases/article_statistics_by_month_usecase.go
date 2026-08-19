package statisticsusecases

import (
	"errors"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ArticleStatisticsByMonthUseCase struct {
	statisticsRepository repository.StatisticsRepository
}

// @inject
func NewArticleStatisticsByMonthUseCase(statisticsRepository repository.StatisticsRepository) *ArticleStatisticsByMonthUseCase {
	return &ArticleStatisticsByMonthUseCase{
		statisticsRepository: statisticsRepository,
	}
}

func (this *ArticleStatisticsByMonthUseCase) Execute(year, month int, journalID *uint, publisherID *uint) (*entity.ArticleStatisticsByMonthEntity, error) {
	if year == 0 && month == 0 {
		stats, err := this.statisticsRepository.GetArticleStatisticsByLast30Days(journalID, publisherID)
		if err != nil {
			return nil, err
		}
		return stats, nil
	}

	if month < 1 || month > 12 {
		return nil, errors.New("invalid month parameter. Must be between 1 and 12")
	}
	if year < 2000 || year > 2100 {
		return nil, errors.New("invalid year parameter. Must be between 2000 and 2100")
	}

	stats, err := this.statisticsRepository.GetArticleStatisticsByMonth(year, month, journalID, publisherID)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
