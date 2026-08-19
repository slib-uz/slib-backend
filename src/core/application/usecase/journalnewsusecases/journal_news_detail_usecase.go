package journalnewsusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalNewsDetailUseCase struct {
	repository repository.JournalNewsRepository
}

// @inject
func NewJournalNewsDetailUseCase(repository repository.JournalNewsRepository) *JournalNewsDetailUseCase {
	return &JournalNewsDetailUseCase{repository: repository}
}

func (this *JournalNewsDetailUseCase) Execute(id uint) (*entity.JournalNewsEntity, error) {
	return this.repository.GetByID(id)
}
