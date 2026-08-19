package spellcheckusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type SpellCheckResultStatsByPublisherUseCase struct {
	repository repository.SpellCheckResultRepository
}

// @inject
func NewSpellCheckResultStatsByPublisherUseCase(repository repository.SpellCheckResultRepository) *SpellCheckResultStatsByPublisherUseCase {
	return &SpellCheckResultStatsByPublisherUseCase{repository: repository}
}

func (this *SpellCheckResultStatsByPublisherUseCase) Execute() ([]*entity.PublisherSpellcheckStatsEntity, error) {
	stats, err := this.repository.StatsByPublisher()
	if err != nil {
		return nil, err
	}

	return stats, nil
}
