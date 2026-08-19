package support_dialog_usecase

import (
	"fmt"

	service2 "slib.uz/src/core/application/service"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type SupportDialogCreateAnswerUseCase struct {
	repository              repository.SupportDialogRepository
	sendNotificationService *service2.SendNotificationService
	sendEmailService        *service2.SendEmailService
}

// @inject
func NewSupportDialogCreateAnswerUseCase(
	repository repository.SupportDialogRepository,
	sendNotificationService *service2.SendNotificationService,
	sendEmailService *service2.SendEmailService,
) *SupportDialogCreateAnswerUseCase {
	return &SupportDialogCreateAnswerUseCase{
		repository:              repository,
		sendNotificationService: sendNotificationService,
		sendEmailService:        sendEmailService,
	}
}

func (this SupportDialogCreateAnswerUseCase) Execute(userID uint, supportDialog *entity.SupportDialogEntity) error {
	supportDialog.OwnerID = userID
	supportDialog.MessageType = enum.Answer
	if err := this.sendNotificationService.NewSupportAnswer([]uint{supportDialog.ChatID}, supportDialog.ChatID, supportDialog.Message); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	if err := this.sendEmailService.NewSupportAnswer([]uint{supportDialog.ChatID}, supportDialog.ChatID, supportDialog.Message); err != nil {
		return err
	}
	return this.repository.CreateAnswer(supportDialog)
}
