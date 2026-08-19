package articleusecases

import (
	"slib.uz/src/core/application/mapper"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type EditionArticlesListUseCase struct {
	repository repository.PublishedArticleRepository
}

// @inject
func NewEditionArticlesListUseCase(repository repository.PublishedArticleRepository) *EditionArticlesListUseCase {
	return &EditionArticlesListUseCase{repository: repository}
}

func (this *EditionArticlesListUseCase) Execute(journalID, editionID uint, unassigned bool, page, pageSize int) (*entity.PagingEntity[entity.ArticleInputEntity], error) {
	paging, err := this.repository.GetByEdition(journalID, editionID, unassigned, page, pageSize)
	if err != nil {
		return nil, err
	}
	return entity.NewPagingEntity(page, pageSize, paging.Total, mapper.ArticleEntityListToInput(paging.Items)), nil
}
