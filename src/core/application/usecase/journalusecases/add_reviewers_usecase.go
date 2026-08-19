package journalusecases

import "slib.uz/src/core/domain/ports/repository"

type AddReviewersUsecase struct {
	repository repository.JournalRepository
}

// @inject
func NewAddReviewersUsecase(repository repository.JournalRepository) *AddReviewersUsecase {
	return &AddReviewersUsecase{repository: repository}
}

func (this *AddReviewersUsecase) Execute(journalID uint, reviewerIds []uint) error {
	return this.repository.AddReviewers(journalID, reviewerIds)
}
