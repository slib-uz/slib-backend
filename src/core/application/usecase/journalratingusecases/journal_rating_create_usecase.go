package journalratingusecases

import (
	"errors"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
)

type JournalRatingCreateUseCase struct {
	repository repository.JournalRatingRepository
}

// @inject
func NewJournalRatingCreateUseCase(repository repository.JournalRatingRepository) *JournalRatingCreateUseCase {
	return &JournalRatingCreateUseCase{repository: repository}
}

func (this *JournalRatingCreateUseCase) Execute(userID uint, rating *entity.JournalRatingEntity) error {
	rating.UserID = userID
	if rating.Stars < 1 || rating.Stars > 5 {
		return infraError.Wrap(errors.New("stars must be between 1 and 5"))
	}
	return this.repository.Create(rating)
}
