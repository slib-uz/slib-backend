package publisherusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type PublisherDetailUseCase struct {
	repository repository.PublisherRepository
}

// @inject
func NewPublisherDetailUseCase(repository repository.PublisherRepository) *PublisherDetailUseCase {
	return &PublisherDetailUseCase{repository: repository}
}

func (this PublisherDetailUseCase) Execute(id uint) (*entity.PublisherEntity, error) {
	ent, err := this.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return ent, nil

}
