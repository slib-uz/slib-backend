package ethicspolicyusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type EthicsPolicyListUseCase struct {
	repository repository.EthicsPolicyRepository
}

// @inject
func NewEthicsPolicyListUseCase(repository repository.EthicsPolicyRepository) *EthicsPolicyListUseCase {
	return &EthicsPolicyListUseCase{repository: repository}
}

func (this *EthicsPolicyListUseCase) Execute(page, pageSize int) (*entity.PagingEntity[entity.EthicsPolicyEntity], error) {
	return this.repository.GetByPaging(page, pageSize)
}
