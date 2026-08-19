package statisticsusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type SiteStatisticsUseCase struct {
	statisticsRepository repository.StatisticsRepository
}

// @inject
func NewSiteStatisticsUseCase(statisticsRepository repository.StatisticsRepository) *SiteStatisticsUseCase {
	return &SiteStatisticsUseCase{
		statisticsRepository: statisticsRepository,
	}
}

func (this *SiteStatisticsUseCase) Execute(journalID *uint, publisherID *uint, institutionID *uint) (*entity.SiteStatisticsEntity, error) {
	stats, err := this.statisticsRepository.GetSiteStatistics(journalID, publisherID, institutionID)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
