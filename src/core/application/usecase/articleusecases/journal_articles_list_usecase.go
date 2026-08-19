package articleusecases

import (
	"slib.uz/src/core/application/mapper"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalArticlesListUseCase struct {
	repository repository.PublishedArticleRepository
}

// @inject
func NewJournalArticlesListUseCase(repository repository.PublishedArticleRepository) *JournalArticlesListUseCase {
	return &JournalArticlesListUseCase{repository: repository}
}

func (this *JournalArticlesListUseCase) Execute(journalID uint, page, pageSize int, search string, authorSearch, tag string, studyFieldIds []uint, year int) (*entity2.PagingEntity[entity2.ArticleInputEntity], error) {
	paging, err := this.repository.GetByJournal(journalID, page, pageSize, search, authorSearch, tag, studyFieldIds, year)
	if err != nil {
		return nil, err
	}
	return entity2.NewPagingEntity(page, pageSize, paging.Total, mapper.ArticleEntityListToInput(paging.Items)), nil

}
