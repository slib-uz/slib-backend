package institutionusecases

import (
	"slib.uz/src/core/domain/ports/repository"
)

type InstitutionDeleteUseCase struct {
	repo repository.InstitutionRepository
}

// @inject
func NewInstitutionDeleteUseCase(repo repository.InstitutionRepository) *InstitutionDeleteUseCase {
	return &InstitutionDeleteUseCase{repo: repo}
}

func (this *InstitutionDeleteUseCase) Execute(id uint) error {
	return this.repo.Delete(id)
}
