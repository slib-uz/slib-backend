package newsusecases

import (
	"fmt"

	"slib.uz/src/core/domain/ports/cache"
	"slib.uz/src/core/domain/ports/repository"
)

type UpdateNewsViewsCountUseCase struct {
	repository repository.NewsRepository
	cache      cache.NewsViewsCountCache
}

// @inject
func NewUpdateNewsViewsCountUseCase(repository repository.NewsRepository, cache cache.NewsViewsCountCache) *UpdateNewsViewsCountUseCase {
	return &UpdateNewsViewsCountUseCase{repository: repository, cache: cache}
}

func (uc *UpdateNewsViewsCountUseCase) Execute() error {
	viewCounts, err := uc.cache.GetAll()
	if err != nil {
		return fmt.Errorf("get news views count from cache: %w", err)
	}
	return uc.repository.UpdateViewsCount(viewCounts)
}
