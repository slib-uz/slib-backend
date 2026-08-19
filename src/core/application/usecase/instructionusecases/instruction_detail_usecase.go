package instructionusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type InstructionDetailUseCase struct {
	repository repository.InstructionRepository
}

// @inject
func NewInstructionDetailUseCase(repository repository.InstructionRepository) *InstructionDetailUseCase {
	return &InstructionDetailUseCase{repository: repository}
}

func (this *InstructionDetailUseCase) Execute(key string) (*entity.InstructionEntity, error) {
	return this.repository.GetByKey(key)
}
