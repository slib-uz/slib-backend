package journalusecases

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ChangeJournalStatusUseCase struct {
	repo repository.JournalRepository
}

// @inject
func NewChangeJournalStatusUseCase(repo repository.JournalRepository) *ChangeJournalStatusUseCase {
	return &ChangeJournalStatusUseCase{
		repo: repo,
	}
}

func (this *ChangeJournalStatusUseCase) Execute(user *entity.UserBasicEntity, journalID uint, isActive bool) error {
	if !user.IsAdmin {
		return response.PermissionDeniedError
	}

	return this.repo.UpdateStatus(journalID, isActive)
}
