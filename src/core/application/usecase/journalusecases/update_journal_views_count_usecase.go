package journalusecases

import (
	"fmt"

	"slib.uz/src/core/domain/ports/cache"
	"slib.uz/src/core/domain/ports/repository"
)

type UpdateJournalViewsCountUseCase struct {
	repository repository.JournalRepository
	cache      cache.JournalViewsCountCache
}

// @inject
func NewUpdateJournalViewsCountUseCase(repository repository.JournalRepository, cache cache.JournalViewsCountCache) *UpdateJournalViewsCountUseCase {
	return &UpdateJournalViewsCountUseCase{repository: repository, cache: cache}
}

func (uc *UpdateJournalViewsCountUseCase) Execute() error {
	viewCounts, err := uc.cache.GetAll()
	if err != nil {
		return fmt.Errorf("get journal views count from cache: %w", err)
	}
	return uc.repository.UpdateViewsCount(viewCounts)
}
