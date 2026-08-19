package newsusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/cache"
	"slib.uz/src/core/domain/ports/repository"
)

type NewsRetrieveUseCase struct {
	repository repository.NewsRepository
	viewsCache cache.NewsViewsCountCache
}

// @inject
func NewNewsRetrieveUseCase(repository repository.NewsRepository, viewsCache cache.NewsViewsCountCache) *NewsRetrieveUseCase {
	return &NewsRetrieveUseCase{repository: repository, viewsCache: viewsCache}
}

func (this *NewsRetrieveUseCase) Execute(newsID uint, userKey string) (*entity.NewsEntity, error) {
	news, err := this.repository.GetByID(newsID)
	if err != nil {
		return nil, err
	}

	count, err := this.viewsCache.Add(userKey, newsID)
	if err == nil {
		news.ViewsCount = news.ViewsCount + uint(count)
	}

	return news, nil
}
