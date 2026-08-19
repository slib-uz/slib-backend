package institutionusecases

import (
	"slib.uz/src/core/domain/ports/repository"
)

type InstitutionSetPublishersUseCase struct {
	repo repository.InstitutionRepository
}

// @inject
func NewInstitutionSetPublishersUseCase(repo repository.InstitutionRepository) *InstitutionSetPublishersUseCase {
	return &InstitutionSetPublishersUseCase{repo: repo}
}

func (this *InstitutionSetPublishersUseCase) Execute(institutionID uint, publisherIDs []uint) error {
	return this.repo.SetPublishers(institutionID, publisherIDs)
}
