package articleusecases

import (
	"slib.uz/src/core/application/mapper"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type TopArticlesUseCase struct {
	repository repository.PublishedArticleRepository
}

// @inject
func NewTopArticlesUseCase(repository repository.PublishedArticleRepository) *TopArticlesUseCase {
	return &TopArticlesUseCase{repository: repository}
}

func (this *TopArticlesUseCase) Execute(page, pageSize int) (*entity2.PagingEntity[entity2.ArticleInputEntity], error) {

	paging, err := this.repository.GetTopArticles(page, pageSize)

	if err != nil {
		return nil, err
	}

	return entity2.NewPagingEntity(page, pageSize, paging.Total, mapper.ArticleEntityListToInput(paging.Items)), nil

}
