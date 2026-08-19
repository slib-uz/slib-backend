package journal_applications_usecases

import (
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type ReviewApplicationUseCase struct {
	repository repository.JournalApplicationRepository
}

// @inject
func NewReviewApplicationUseCase(repository repository.JournalApplicationRepository) *ReviewApplicationUseCase {
	return &ReviewApplicationUseCase{repository: repository}
}

func (this *ReviewApplicationUseCase) Execute(id uint, status enum.Status, rejectionReason *string) error {
	return this.repository.SetStatus(id, enum.Status(status), rejectionReason)
}
