package newsusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type NewsListUseCase struct {
	repository repository.NewsRepository
}

// @inject
func NewNewsListUseCase(repository repository.NewsRepository) *NewsListUseCase {
	return &NewsListUseCase{repository: repository}
}

func (this NewsListUseCase) Execute(page, pageSize int, ordering string, categoryID uint) (*entity2.PagingEntity[entity2.NewsEntity], error) {
	paging, err := this.repository.GetByPaging(page, pageSize, ordering, categoryID)

	if err != nil {
		return nil, err
	}

	return paging, nil
}
