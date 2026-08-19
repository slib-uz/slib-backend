package articleusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ArticleAndJournalCountUseCase struct {
	articleRepository repository.ArticleRepository
	journalRepository repository.JournalRepository
}

// @inject
func NewArticleAndJournalCountUseCase(articleRepository repository.ArticleRepository, journalRepository repository.JournalRepository) *ArticleAndJournalCountUseCase {
	return &ArticleAndJournalCountUseCase{articleRepository: articleRepository, journalRepository: journalRepository}
}

func (this *ArticleAndJournalCountUseCase) Execute() (*entity.CountStatsEntity, error) {
	journalCount, err := this.journalRepository.GetJournalCount()
	if err != nil {
		return nil, err
	}
	articleCount, err := this.articleRepository.GetArticleCount()
	if err != nil {
		return nil, err
	}
	authorCount, err := this.articleRepository.GetAuthorsCount()
	if err != nil {
		return nil, err
	}
	return entity.NewCountStatsEntity(journalCount, articleCount, authorCount), nil
}
