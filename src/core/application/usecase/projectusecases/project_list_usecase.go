package projectusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ProjectListUseCase struct {
	repository repository.ProjectRepository
}

// @inject
func NewProjectListUseCase(repository repository.ProjectRepository) *ProjectListUseCase {
	return &ProjectListUseCase{repository: repository}
}

func (this *ProjectListUseCase) Execute(page, pageSize int) (*entity2.PagingEntity[entity2.ProjectEntity], error) {
	paging, err := this.repository.GetAll(page, pageSize)

	if err != nil {
		return nil, err
	}

	return paging, nil
}
