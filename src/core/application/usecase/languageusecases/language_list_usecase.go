package languageusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type LanguageListUseCase struct {
	repository repository.LanguageRepository
}

// @inject
func NewLanguageListUseCase(repository repository.LanguageRepository) *LanguageListUseCase {
	return &LanguageListUseCase{repository: repository}
}

func (this *LanguageListUseCase) Execute() ([]*entity.LanguageEntity, error) {

	languages, err := this.repository.GetAll()
	if err != nil {
		return nil, err
	}

	return languages, nil
}
