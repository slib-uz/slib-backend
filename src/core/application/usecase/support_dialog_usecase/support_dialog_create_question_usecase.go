package support_dialog_usecase

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type SupportDialogCreateQuestionUseCase struct {
	repository repository.SupportDialogRepository
}

// @inject
func NewSupportDialogCreateQuestionUseCase(repository repository.SupportDialogRepository) *SupportDialogCreateQuestionUseCase {
	return &SupportDialogCreateQuestionUseCase{repository: repository}
}

func (this SupportDialogCreateQuestionUseCase) Execute(userID uint, supportDialog *entity.SupportDialogEntity) error {
	supportDialog.OwnerID = userID
	supportDialog.MessageType = enum.Question
	supportDialog.ChatID = userID
	return this.repository.CreateQuestion(supportDialog)
}
