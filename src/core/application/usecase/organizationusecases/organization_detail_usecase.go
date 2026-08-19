package organizationusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type OrganizationDetailUseCase struct {
	repository repository.OrganizationRepository
}

// @inject
func NewOrganizationDetailUseCase(repository repository.OrganizationRepository) *OrganizationDetailUseCase {
	return &OrganizationDetailUseCase{repository: repository}
}

func (this *OrganizationDetailUseCase) Execute(id uint) (*entity.OrganizationEntity, error) {
	return this.repository.GetByID(id)
}
