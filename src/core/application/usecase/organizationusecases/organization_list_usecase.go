package organizationusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type OrganizationListUseCase struct {
	repository repository.OrganizationRepository
}

// @inject
func NewOrganizationListUseCase(repository repository.OrganizationRepository) *OrganizationListUseCase {
	return &OrganizationListUseCase{repository: repository}
}

func (this *OrganizationListUseCase) Execute(page, size int, tin, name, address *string) (*entity2.PagingEntity[entity2.OrganizationEntity], error) {
	entPaging, err := this.repository.GetList(page, size, tin, name, address)
	if err != nil {
		return nil, err
	}

	return entPaging, nil
}
