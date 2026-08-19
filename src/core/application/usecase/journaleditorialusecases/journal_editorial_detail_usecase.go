package journaleditorialusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalEditorialDetailUseCase struct {
	repository repository.JournalEditorialRepository
}

// @inject
func NewJournalEditorialDetailUseCase(repository repository.JournalEditorialRepository) *JournalEditorialDetailUseCase {
	return &JournalEditorialDetailUseCase{repository: repository}
}

func (this *JournalEditorialDetailUseCase) Execute(id uint) (*entity.JournalEditorialEntity, error) {
	return this.repository.GetByID(id)
}
