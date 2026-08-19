package institutionusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type InstitutionUpdateUseCase struct {
	repo repository.InstitutionRepository
}

// @inject
func NewInstitutionUpdateUseCase(repo repository.InstitutionRepository) *InstitutionUpdateUseCase {
	return &InstitutionUpdateUseCase{repo: repo}
}

func (this *InstitutionUpdateUseCase) Execute(id uint, entity *entity2.InstitutionEntity) error {
	return this.repo.Update(id, entity)
}
