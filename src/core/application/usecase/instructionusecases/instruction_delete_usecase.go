package instructionusecases

import "slib.uz/src/core/domain/ports/repository"

type InstructionDeleteUseCase struct {
	repository repository.InstructionRepository
}

// @inject
func NewInstructionDeleteUseCase(repository repository.InstructionRepository) *InstructionDeleteUseCase {
	return &InstructionDeleteUseCase{repository: repository}
}

func (this *InstructionDeleteUseCase) Execute(id uint) error {
	return this.repository.Delete(id)
}
