package publisherusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type PublisherListUseCase struct {
	repository repository.PublisherRepository
}

// @inject
func NewPublisherListUseCase(repository repository.PublisherRepository) *PublisherListUseCase {
	return &PublisherListUseCase{repository: repository}
}

func (this PublisherListUseCase) Execute(page, size int, name, tin string, institutionID uint, unassigned bool) (*entity2.PagingEntity[entity2.PublisherEntity], error) {

	entPaging, err := this.repository.GetListByPage(page, size, name, tin, institutionID, unassigned)

	if err != nil {
		return nil, err
	}

	return entPaging, nil
}
