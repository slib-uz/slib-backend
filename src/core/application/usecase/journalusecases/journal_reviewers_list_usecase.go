package journalusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalReviewersListUsecase struct {
	repository repository.ReviewerRepository
}

// @inject
func NewJournalReviewersListUsecase(repository repository.ReviewerRepository) *JournalReviewersListUsecase {
	return &JournalReviewersListUsecase{repository: repository}
}

func (this *JournalReviewersListUsecase) Execute(journalID uint) ([]*entity.ReviewerEntity, error) {
	reviewers, err := this.repository.GetByJournalID(journalID)
	if err != nil {
		return nil, err
	}

	return reviewers, nil
}
