package soatousecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type SoatoListUseCase struct {
	repo repository.SoatoRepository
}

// @inject
func NewSoatoListUseCase(repo repository.SoatoRepository) *SoatoListUseCase {
	return &SoatoListUseCase{repo: repo}
}

func (this *SoatoListUseCase) Execute(page, size int, isChildren bool) (*entity.PagingEntity[entity.SoatoClassifierEntity], error) {
	return this.repo.GetAll(page, size, isChildren)
}
