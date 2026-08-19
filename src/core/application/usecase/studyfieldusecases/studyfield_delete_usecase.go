package studyfieldusecases

import "slib.uz/src/core/domain/ports/repository"

type StudyFieldDeleteUseCase struct {
	repository repository.StudyFieldRepository
}

// @inject
func NewStudyFieldDeleteUseCase(repository repository.StudyFieldRepository) *StudyFieldDeleteUseCase {
	return &StudyFieldDeleteUseCase{repository: repository}
}

func (this *StudyFieldDeleteUseCase) Execute(id uint) error {
	return this.repository.Delete(id)
}
