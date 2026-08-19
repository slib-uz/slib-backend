package articleusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ArticleDownloadUseCase struct {
	articleRepository     repository.ArticleRepository
	articleFileRepository repository.ArticleFileRepository
}

// @inject
func NewArticleDownloadUseCase(articleRepository repository.ArticleRepository, articleFileRepository repository.ArticleFileRepository) *ArticleDownloadUseCase {
	return &ArticleDownloadUseCase{articleRepository: articleRepository, articleFileRepository: articleFileRepository}
}

func (this *ArticleDownloadUseCase) Execute(articleID uint) (*entity.ArticleDownloadEntity, error) {
	article, err := this.articleRepository.FindByID(articleID)
	if err != nil {
		return nil, err
	}

	articleFile, aErr := this.articleFileRepository.GetFile(article.ContentFile)
	if aErr != nil {
		return nil, aErr
	}

	return entity.NewArticleDownloadEntity(articleFile.PresignedURL, articleFile.BinaryFile), nil
}
