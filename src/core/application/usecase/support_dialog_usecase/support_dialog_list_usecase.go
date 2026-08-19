package support_dialog_usecase

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type SupportDialogQuestionListUseCase struct {
	repository repository.SupportDialogRepository
}

// @inject
func NewSupportDialogQuestionListUseCase(repository repository.SupportDialogRepository) *SupportDialogQuestionListUseCase {
	return &SupportDialogQuestionListUseCase{repository: repository}
}

func (this SupportDialogQuestionListUseCase) Execute(page, pageSize int, ordering string, userID uint) (*entity2.PagingEntity[entity2.SupportDialogEntity], error) {
	paging, err := this.repository.GetByPaging(page, pageSize, ordering, userID)

	if err != nil {
		return nil, err
	}

	return paging, nil
}
