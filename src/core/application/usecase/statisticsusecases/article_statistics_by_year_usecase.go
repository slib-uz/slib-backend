package statisticsusecases

import (
	"errors"
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ArticleStatisticsByYearUseCase struct {
	statisticsRepository repository.StatisticsRepository
}

// @inject
func NewArticleStatisticsByYearUseCase(statisticsRepository repository.StatisticsRepository) *ArticleStatisticsByYearUseCase {
	return &ArticleStatisticsByYearUseCase{
		statisticsRepository: statisticsRepository,
	}
}

func (this *ArticleStatisticsByYearUseCase) Execute(year int, journalID *uint, publisherID *uint) (*entity.ArticleStatisticsByYearEntity, error) {
	if year == 0 {
		year = time.Now().Year()
	}

	if year < 2000 || year > 2100 {
		return nil, errors.New("invalid year parameter. Must be between 2000 and 2100")
	}

	stats, err := this.statisticsRepository.GetArticleStatisticsByYear(year, journalID, publisherID)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
