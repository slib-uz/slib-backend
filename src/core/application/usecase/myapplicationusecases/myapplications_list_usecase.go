package myapplicationusecases

import (
	"slib.uz/src/core/application/mapper"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type MyApplicationsListUseCase struct {
	repository repository.ApplicationRepository
}

// @inject
func NewMyApplicationsListUseCase(repository repository.ApplicationRepository) *MyApplicationsListUseCase {
	return &MyApplicationsListUseCase{repository: repository}
}

func (this *MyApplicationsListUseCase) Execute(userId uint, page, pageSize int, journalID uint) (*entity2.PagingEntity[entity2.ApplicationBasicEntity], error) {

	paging, err := this.repository.GetByUserID(userId, page, pageSize, journalID)
	if err != nil {
		return nil, err
	}

	return entity2.NewPagingEntity(
		page,
		pageSize,
		paging.Total,
		mapper.ApplicationEntityListToBasicEntityList(paging.Items),
	), nil

}
