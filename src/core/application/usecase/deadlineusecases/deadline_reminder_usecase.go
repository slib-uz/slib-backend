package deadlineusecases

import (
	"fmt"

	"github.com/labstack/gommon/log"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type DeadlineReminderUseCase struct {
	reviewStageRepo         repository.ReviewStageRepository
	sendNotificationService *service.SendNotificationService
	roleRepo                repository.UserRoleRepository
}

// @inject
func NewDeadlineReminderUseCase(
	reviewStageRepo repository.ReviewStageRepository,
	sendNotificationService *service.SendNotificationService,
	roleRepo repository.UserRoleRepository,
) *DeadlineReminderUseCase {
	return &DeadlineReminderUseCase{
		reviewStageRepo:         reviewStageRepo,
		sendNotificationService: sendNotificationService,
		roleRepo:                roleRepo,
	}
}

func (uc *DeadlineReminderUseCase) Execute() error {
	reviewStages, err := uc.reviewStageRepo.GetExpiringStages()
	if err != nil {
		return fmt.Errorf("get expiring stages(DeadlineReminderUseCase): %w", err)
	}

	for _, s := range reviewStages {
		if s.Stage != enum.PaymentStage {
			uc.sendNotificationToJournalMembers(s)
		} else {
			uc.sendNotificationToAuthor(s)
		}
	}

	return nil
}

func (uc *DeadlineReminderUseCase) sendNotificationToAuthor(reviewStage *entity.ReviewStageEntity) {
	app := reviewStage.Application

	if err := uc.sendNotificationService.ReviewDeadlineRemind(
		reviewStage.Stage,
		[]uint{app.UserID},
		app.ID,
		app.Number,
		reviewStage.Deadline,
	); err != nil {
		log.Errorf("send deadline reminder notification to applicant(DeadlineReminderUseCase): %v", err)
	}
}

func (uc *DeadlineReminderUseCase) sendNotificationToJournalMembers(reviewStage *entity.ReviewStageEntity) {
	journalMemberIds, err := uc.roleRepo.GetJournalMemberIds(reviewStage.Application.JournalID)
	if err != nil {
		log.Errorf("get journal member ids(DeadlineReminderUseCase): %v", err)
		return
	}

	if len(journalMemberIds) == 0 {
		log.Warnf("no journal members found for journal id %d", reviewStage.Application.JournalID)
		return
	}

	if err := uc.sendNotificationService.ReviewDeadlineRemind(
		reviewStage.Stage,
		journalMemberIds,
		reviewStage.ApplicationID,
		reviewStage.Application.Number,
		reviewStage.Deadline,
	); err != nil {
		log.Errorf("send deadline reminder notification(DeadlineReminderUseCase): %v", err)
	}
}
