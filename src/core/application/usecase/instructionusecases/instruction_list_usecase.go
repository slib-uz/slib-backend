package instructionusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type InstructionListUseCase struct {
	repository repository.InstructionRepository
}

// @inject
func NewInstructionListUseCase(repository repository.InstructionRepository) *InstructionListUseCase {
	return &InstructionListUseCase{repository: repository}
}

func (this *InstructionListUseCase) Execute() ([]*entity.InstructionEntity, error) {
	return this.repository.GetAll()
}
