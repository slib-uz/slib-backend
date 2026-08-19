package studyfieldusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type StudyFieldCreateUseCase struct {
	repository repository.StudyFieldRepository
}

// @inject
func NewStudyFieldCreateUseCase(repository repository.StudyFieldRepository) *StudyFieldCreateUseCase {
	return &StudyFieldCreateUseCase{repository: repository}
}

func (this *StudyFieldCreateUseCase) Execute(studyField *entity.StudyFieldEntity) error {
	return this.repository.Create(studyField)
}
