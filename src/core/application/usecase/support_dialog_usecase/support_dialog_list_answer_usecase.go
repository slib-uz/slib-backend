package support_dialog_usecase

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type SupportDialogAnswerListUseCase struct {
	repository repository.SupportDialogRepository
}

// @inject
func NewSupportDialogAnswerListUseCase(repository repository.SupportDialogRepository) *SupportDialogAnswerListUseCase {
	return &SupportDialogAnswerListUseCase{repository: repository}
}

func (this SupportDialogAnswerListUseCase) Execute(page, pageSize int, ordering string, chatID uint) (*entity2.PagingEntity[entity2.SupportDialogEntity], error) {
	if err := this.repository.MarkAsRead(chatID); err != nil {
		return nil, err
	}

	paging, err := this.repository.GetByChatID(page, pageSize, ordering, chatID)

	if err != nil {
		return nil, err
	}

	return paging, nil
}
