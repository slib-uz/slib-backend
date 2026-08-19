package chiefeditorusecase

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ChiefEditorByJournalIDUseCase struct {
	repository repository.ChiefEditorRepository
}

// @inject
func NewChiefEditorByJournalIDUseCase(repository repository.ChiefEditorRepository) *ChiefEditorByJournalIDUseCase {
	return &ChiefEditorByJournalIDUseCase{repository: repository}
}

func (this *ChiefEditorByJournalIDUseCase) Execute(journalID uint) ([]*entity.ChiefEditorEntity, error) {
	chiefEditors, err := this.repository.GetByJournalID(journalID)
	if err != nil {
		return nil, err
	}

	return chiefEditors, nil
}
