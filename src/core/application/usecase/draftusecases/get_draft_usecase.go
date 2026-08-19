package draftusecases

import (
	"errors"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type GetDraftUseCase struct {
	repository repository.DraftRepository
}

// @inject
func NewGetDraftUseCase(repository repository.DraftRepository) *GetDraftUseCase {
	return &GetDraftUseCase{repository: repository}
}

func (this *GetDraftUseCase) Execute(key string) (*entity.DraftEntity, error) {
	draft, err := this.repository.GetByKey(key)
	if err != nil && !errors.Is(err, response.NotFoundError) {
		return nil, err
	}
	return draft, nil
}
