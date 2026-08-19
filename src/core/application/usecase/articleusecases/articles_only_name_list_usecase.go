package articleusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ArticlesOnlyNameListUseCase struct {
	legacyAuthorRepository repository.LegacyAuthorRepository
	articleRepository      repository.ArticleRepository
}

// @inject
func NewArticlesOnlyNameListUseCase(legacyAuthorRepository repository.LegacyAuthorRepository, articleRepository repository.ArticleRepository) *ArticlesOnlyNameListUseCase {
	return &ArticlesOnlyNameListUseCase{legacyAuthorRepository: legacyAuthorRepository, articleRepository: articleRepository}
}

func (this *ArticlesOnlyNameListUseCase) Execute(fullName string, page, pageSize int) (*entity.PagingEntity[entity.ArticleOnlyNameEntity], error) {
	ids, err := this.legacyAuthorRepository.GetIDsByFullName(fullName)
	if err != nil {
		return nil, err
	}

	return this.articleRepository.OnlyArticleNameByLegacyAuthorIDs(ids, page, pageSize)
}
