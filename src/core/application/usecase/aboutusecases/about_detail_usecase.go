package aboutusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type AboutDetailUseCase struct {
	repo repository.AboutRepository
}

// @inject
func NewAboutDetailUseCase(repo repository.AboutRepository) *AboutDetailUseCase {
	return &AboutDetailUseCase{repo: repo}
}

func (this *AboutDetailUseCase) Execute(id uint) (*entity.AboutEntity, error) {
	return this.repo.GetByID(id)
}
