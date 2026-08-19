package antiplagusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type AntiPlagStatsByPublisherUseCase struct {
	repository repository.AntiPlagRepository
}

// @inject
func NewAntiPlagStatsByPublisherUseCase(repository repository.AntiPlagRepository) *AntiPlagStatsByPublisherUseCase {
	return &AntiPlagStatsByPublisherUseCase{repository: repository}
}

func (this *AntiPlagStatsByPublisherUseCase) Execute() ([]*entity.PublisherAntiplagStatsEntity, error) {
	stats, err := this.repository.StatsByPublisher()
	if err != nil {
		return nil, err
	}
	return stats, nil
}
