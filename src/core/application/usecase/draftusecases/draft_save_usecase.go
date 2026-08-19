package draftusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type DraftSaveUseCase struct {
	repository repository.DraftRepository
}

// @inject
func NewDraftSaveUseCase(repository repository.DraftRepository) *DraftSaveUseCase {
	return &DraftSaveUseCase{repository: repository}
}

func (this *DraftSaveUseCase) Execute(draft *entity.DraftEntity) error {
	return this.repository.Save(draft)
}
