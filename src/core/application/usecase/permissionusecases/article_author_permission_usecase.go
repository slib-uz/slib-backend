package permissionusecases

import (
	"slib.uz/src/core/domain/ports/repository"
)

type ArticleAuthorPermissionUseCase struct {
	authorRepository  repository.AuthorRepository
	articleRepository repository.ArticleRepository
}

// @inject
func NewArticleAuthorPermissionUseCase(authorRepository repository.AuthorRepository, articleRepository repository.ArticleRepository) *ArticleAuthorPermissionUseCase {
	return &ArticleAuthorPermissionUseCase{authorRepository: authorRepository, articleRepository: articleRepository}
}

func (this *ArticleAuthorPermissionUseCase) Execute(userScienceID string, articleID uint) (bool, error) {
	article, err := this.articleRepository.FindByIDWithAuthors(articleID)
	if err != nil {
		return false, err
	}

	for _, author := range article.CoAuthors {
		if author.ScienceID == userScienceID {
			return true, nil
		}
	}
	return false, nil
}
