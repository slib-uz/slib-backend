package projectusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ProjectDetailUseCase struct {
	repository repository.ProjectRepository
}

// @inject
func NewProjectDetailUseCase(repository repository.ProjectRepository) *ProjectDetailUseCase {
	return &ProjectDetailUseCase{repository: repository}
}

func (this *ProjectDetailUseCase) Execute(id uint) (*entity.ProjectEntity, error) {
	ent, err := this.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return ent, nil
}
