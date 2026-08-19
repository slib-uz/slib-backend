package journalratingusecases

import (
	"slib.uz/src/core/domain/ports/repository"
)

type JournalRatingDeleteUseCase struct {
	repository repository.JournalRatingRepository
}

// @inject
func NewJournalRatingDeleteUseCase(repository repository.JournalRatingRepository) *JournalRatingDeleteUseCase {
	return &JournalRatingDeleteUseCase{repository: repository}
}

func (this *JournalRatingDeleteUseCase) Execute(userID, id uint) error {
	return this.repository.DeleteByIDAndUserID(id, userID)
}
