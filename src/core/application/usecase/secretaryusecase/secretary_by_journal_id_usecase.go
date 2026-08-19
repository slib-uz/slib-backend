package secretaryusecase

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type SecretaryByJournalIDUseCase struct {
	repository repository.SecretaryRepository
}

// @inject
func NewSecretaryByJournalIDUseCase(repository repository.SecretaryRepository) *SecretaryByJournalIDUseCase {
	return &SecretaryByJournalIDUseCase{repository: repository}
}

func (this *SecretaryByJournalIDUseCase) Execute(journalID uint) ([]*entity.SecretaryEntity, error) {
	secretaries, err := this.repository.GetByJournalID(journalID)
	if err != nil {
		return nil, err
	}
	return secretaries, nil
}
