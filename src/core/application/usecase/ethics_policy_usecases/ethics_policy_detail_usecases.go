package ethicspolicyusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type EthicsPolicyDetailUseCase struct {
	repository repository.EthicsPolicyRepository
}

// @inject
func NewEthicsPolicyDetailUseCase(repository repository.EthicsPolicyRepository) *EthicsPolicyDetailUseCase {
	return &EthicsPolicyDetailUseCase{repository: repository}
}

func (this *EthicsPolicyDetailUseCase) Execute(id uint) (*entity.EthicsPolicyEntity, error) {
	return this.repository.GetByID(id)
}
