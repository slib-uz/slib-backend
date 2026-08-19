package institutionusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type InstitutionListUseCase struct {
	repository repository.InstitutionRepository
}

// @inject
func NewInstitutionListUseCase(repository repository.InstitutionRepository) *InstitutionListUseCase {
	return &InstitutionListUseCase{repository: repository}
}

func (this *InstitutionListUseCase) Execute(page, size int, tin, name *string) (*entity2.PagingEntity[entity2.InstitutionEntity], error) {
	return this.repository.GetList(page, size, tin, name)
}
