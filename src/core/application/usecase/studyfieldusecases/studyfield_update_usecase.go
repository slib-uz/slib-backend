package studyfieldusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type StudyFieldUpdateUseCase struct {
	repository repository.StudyFieldRepository
}

// @inject
func NewStudyFieldUpdateUseCase(repository repository.StudyFieldRepository) *StudyFieldUpdateUseCase {
	return &StudyFieldUpdateUseCase{repository: repository}
}

func (this *StudyFieldUpdateUseCase) Execute(studyField *entity.StudyFieldEntity) error {
	return this.repository.Update(studyField)
}
