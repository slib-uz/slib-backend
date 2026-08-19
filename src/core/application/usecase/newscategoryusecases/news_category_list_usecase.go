package newscategoryusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type NewsCategoryListUseCase struct {
	repository repository.NewsCategoryRepository
}

// @inject
func NewNewsCategoryListUseCase(repository repository.NewsCategoryRepository) *NewsCategoryListUseCase {
	return &NewsCategoryListUseCase{repository: repository}
}

func (this *NewsCategoryListUseCase) Execute() ([]*entity.NewsCategoryEntity, error) {
	categories, err := this.repository.GetAll()
	if err != nil {
		return nil, err
	}

	return categories, nil
}
