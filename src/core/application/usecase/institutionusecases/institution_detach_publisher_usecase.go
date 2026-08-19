package institutionusecases

import (
	"slib.uz/src/core/domain/ports/repository"
)

type InstitutionDetachPublisherUseCase struct {
	repo repository.InstitutionRepository
}

// @inject
func NewInstitutionDetachPublisherUseCase(repo repository.InstitutionRepository) *InstitutionDetachPublisherUseCase {
	return &InstitutionDetachPublisherUseCase{repo: repo}
}

func (this *InstitutionDetachPublisherUseCase) Execute(institutionID uint, publisherIDs []uint) (int64, error) {
	return this.repo.DetachPublishers(institutionID, publisherIDs)
}
