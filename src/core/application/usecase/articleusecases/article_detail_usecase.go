package articleusecases

import (
	"slib.uz/src/core/application/mapper"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/cache"
	"slib.uz/src/core/domain/ports/repository"
)

type ArticleDetailUseCase struct {
	repository repository.PublishedArticleRepository
	viewsCache cache.ArticleViewsCountCache
}

// @inject
func NewArticleDetailUseCase(repository repository.PublishedArticleRepository, viewsCache cache.ArticleViewsCountCache) *ArticleDetailUseCase {
	return &ArticleDetailUseCase{repository: repository, viewsCache: viewsCache}
}

func (this *ArticleDetailUseCase) Execute(articleID uint, userKey string) (*entity.ArticleInputEntity, error) {
	article, err := this.repository.GetByIDWithRelations(articleID)
	if err != nil {
		return nil, err
	}
	count, err := this.viewsCache.Add(userKey, articleID)
	if err == nil {
		article.ViewsCount = article.ViewsCount + count
	}

	return mapper.ArticleEntityToInput(article), nil
}
