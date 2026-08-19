package journalusecases

import "slib.uz/src/core/domain/ports/repository"

type RemoveReviewersUsecase struct {
	repository repository.ReviewerRepository
}

// @inject
func NewRemoveReviewersUsecase(repository repository.ReviewerRepository) *RemoveReviewersUsecase {
	return &RemoveReviewersUsecase{repository: repository}
}

func (this *RemoveReviewersUsecase) Execute(journalID uint, reviewerID uint) error {

	if err := this.repository.RemoveReviewerFromJournal(journalID, reviewerID); err != nil {
		return err

	}
	return nil
}
