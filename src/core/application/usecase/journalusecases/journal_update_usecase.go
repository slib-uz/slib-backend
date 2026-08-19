package journalusecases

import (
	"slib.uz/src/core/application/validation"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalUpdateUseCase struct {
	repository     repository.JournalRepository
	languageRepo   repository.LanguageRepository
	studyFieldRepo repository.StudyFieldRepository
}

// @inject
func NewJournalUpdateUseCase(
	repository repository.JournalRepository,
	languageRepo repository.LanguageRepository,
	studyFieldRepo repository.StudyFieldRepository,
) *JournalUpdateUseCase {
	return &JournalUpdateUseCase{
		repository:     repository,
		languageRepo:   languageRepo,
		studyFieldRepo: studyFieldRepo,
	}
}

func (this *JournalUpdateUseCase) Execute(journalID uint, e *entity.JournalCreateEntity) error {
	if err := validation.ValidateIDs("Language", e.LanguageIDs, this.languageRepo.ExistingIds); err != nil {
		return err
	}

	if err := validation.ValidateIDs("StudyField", e.StudyFieldIDs, this.studyFieldRepo.ExistingIds); err != nil {
		return err
	}

	return this.repository.UpdateJournal(journalID, e)
}
