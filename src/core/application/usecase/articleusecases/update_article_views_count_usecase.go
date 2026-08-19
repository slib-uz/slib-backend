package articleusecases

import (
	"fmt"

	"slib.uz/src/core/domain/ports/cache"
	"slib.uz/src/core/domain/ports/repository"
)

type UpdateArticleViewsCountUseCase struct {
	repository repository.ArticleRepository
	cache      cache.ArticleViewsCountCache
}

// @inject
func NewUpdateArticleViewsCountUseCase(repository repository.ArticleRepository, cache cache.ArticleViewsCountCache) *UpdateArticleViewsCountUseCase {
	return &UpdateArticleViewsCountUseCase{repository: repository, cache: cache}
}

func (uc *UpdateArticleViewsCountUseCase) Execute() error {
	viewCounts, err := uc.cache.GetAll()
	if err != nil {
		return fmt.Errorf("get article views count from cache: %w", err)
	}
	return uc.repository.UpdateViewsCount(viewCounts)
}
