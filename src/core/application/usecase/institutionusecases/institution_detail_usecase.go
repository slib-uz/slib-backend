package institutionusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type InstitutionDetailUseCase struct {
	repository repository.InstitutionRepository
}

// @inject
func NewInstitutionDetailUseCase(repository repository.InstitutionRepository) *InstitutionDetailUseCase {
	return &InstitutionDetailUseCase{repository: repository}
}

func (this *InstitutionDetailUseCase) Execute(id uint) (*entity.InstitutionEntity, error) {
	return this.repository.GetByID(id)
}
