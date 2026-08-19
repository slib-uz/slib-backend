package institutionusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type InstitutionCreateUseCase struct {
	repo repository.InstitutionRepository
}

// @inject
func NewInstitutionCreateUseCase(repo repository.InstitutionRepository) *InstitutionCreateUseCase {
	return &InstitutionCreateUseCase{repo: repo}
}

func (this *InstitutionCreateUseCase) Execute(entity *entity2.InstitutionEntity) error {
	return this.repo.Create(entity)
}
