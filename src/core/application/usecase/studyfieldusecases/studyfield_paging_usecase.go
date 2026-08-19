package studyfieldusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type StudyFieldPagingUseCase struct {
	repository repository.StudyFieldRepository
}

// @inject
func NewStudyFieldPagingUseCase(repository repository.StudyFieldRepository) *StudyFieldPagingUseCase {
	return &StudyFieldPagingUseCase{repository: repository}
}

func (this *StudyFieldPagingUseCase) Execute(page, pageSize int, search string) (*entity2.PagingEntity[entity2.StudyFieldEntity], error) {
	paging, err := this.repository.GetByPaging(page, pageSize, search)
	if err != nil {
		return nil, err
	}

	return paging, nil
}
