package studyfieldusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type StudyFieldListUseCase struct {
	repository repository.StudyFieldRepository
}

// @inject
func NewStudyFieldListUseCase(repository repository.StudyFieldRepository) *StudyFieldListUseCase {
	return &StudyFieldListUseCase{repository: repository}
}

func (this *StudyFieldListUseCase) Execute(journalID *uint, search string) ([]*entity.StudyFieldEntity, error) {
	var list []*entity.StudyFieldEntity
	var err error

	if journalID != nil {
		list, err = this.repository.GetByJournalID(*journalID)
		if err != nil {
			return nil, err
		}
	} else {
		list, err = this.repository.GetAll(search)
		if err != nil {
			return nil, err
		}
	}

	if journalID == nil {
		return list, nil
	}

	return list, nil
}
