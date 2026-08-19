package support_dialog_usecase

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ChatListUseCase struct {
	repository repository.SupportDialogRepository
}

// @inject
func NewChatListUseCase(repository repository.SupportDialogRepository) *ChatListUseCase {
	return &ChatListUseCase{repository: repository}
}

func (this ChatListUseCase) Execute(page, pageSize int, ordering string) (*entity2.PagingEntity[entity2.ChatEntity], error) {
	paging, err := this.repository.GetChatsByPaging(page, pageSize, ordering)

	if err != nil {
		return nil, err
	}

	return paging, nil
}
