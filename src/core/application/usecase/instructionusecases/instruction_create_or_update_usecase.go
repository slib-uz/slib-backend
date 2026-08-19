package instructionusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type InstructionCreateOrUpdateUseCase struct {
	repository repository.InstructionRepository
}

// @inject
func NewInstructionCreateOrUpdateUseCase(repository repository.InstructionRepository) *InstructionCreateOrUpdateUseCase {
	return &InstructionCreateOrUpdateUseCase{repository: repository}
}

func (this *InstructionCreateOrUpdateUseCase) Execute(instruction *entity.InstructionEntity) error {
	return this.repository.CreateOrUpdate(instruction)
}
