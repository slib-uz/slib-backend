package aboutusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type AboutListUseCase struct {
	repo repository.AboutRepository
}

// @inject
func NewAboutListUseCase(repo repository.AboutRepository) *AboutListUseCase {
	return &AboutListUseCase{repo: repo}
}

func (this *AboutListUseCase) Execute() ([]*entity.AboutEntity, error) {
	list, err := this.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return list, nil
}
