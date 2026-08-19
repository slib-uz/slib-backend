package articleusecases

import (
	"slib.uz/src/core/application/mapper"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type UserArticleListUseCase struct {
	repository repository.PublishedArticleRepository
}

// @inject
func NewUserArticleListUseCase(repository repository.PublishedArticleRepository) *UserArticleListUseCase {
	return &UserArticleListUseCase{repository: repository}
}

func (this *UserArticleListUseCase) Execute(authorScienceID string, page, pageSize int, search string, studyFieldIds []uint, year int, journalID *uint) (*entity2.PagingEntity[entity2.ArticleInputEntity], error) {

	paging, err := this.repository.GetByAuthor(authorScienceID, page, pageSize, search, studyFieldIds, year, journalID)
	if err != nil {
		return nil, err
	}

	return entity2.NewPagingEntity(
			page,
			pageSize,
			paging.Total,
			mapper.ArticleEntityListToInput(paging.Items),
		),
		nil
}
