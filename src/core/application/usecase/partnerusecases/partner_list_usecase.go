package partnerusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type PartnerListUseCase struct {
	repository repository.PartnerRepository
}

// @inject
func NewPartnerListUseCase(repository repository.PartnerRepository) *PartnerListUseCase {
	return &PartnerListUseCase{repository: repository}
}

func (this *PartnerListUseCase) Execute(page, pageSize int) (*entity2.PagingEntity[entity2.PartnerEntity], error) {
	paging, err := this.repository.GetAll(page, pageSize)

	if err != nil {
		return nil, err
	}

	return paging, nil
}
