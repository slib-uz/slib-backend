package aidetectusecases

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/application/tasks"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
)

type AiDetectStatusUpdateUseCase struct {
	repository          repository.AiDetectRepository
	gateway             gateway.AiDetectGateway
	notificationService *service.SendNotificationService
	roleRepo            repository.UserRoleRepository
	outboxRepo          repository.OutboxEventRepository
	outboxTask          *tasks.OutboxEventTaskSender
}

// @inject
func NewAiDetectStatusUpdateUseCase(
	repository repository.AiDetectRepository,
	gateway gateway.AiDetectGateway,
	notificationService *service.SendNotificationService,
	roleRepo repository.UserRoleRepository,
	outboxRepo repository.OutboxEventRepository,
	outboxTask *tasks.OutboxEventTaskSender,
) *AiDetectStatusUpdateUseCase {
	return &AiDetectStatusUpdateUseCase{
		repository:          repository,
		gateway:             gateway,
		notificationService: notificationService,
		roleRepo:            roleRepo,
		outboxRepo:          outboxRepo,
		outboxTask:          outboxTask,
	}
}

func (uc *AiDetectStatusUpdateUseCase) Execute(externalID uint) error {
	result, err := uc.gateway.GetResult(externalID)
	if err != nil {
		return fmt.Errorf("ai detect status update task result error: %w", err)
	}

	id, err := uc.repository.FindIDByExternalID(externalID)
	if err != nil {
		return fmt.Errorf("find id from repository error: %w", err)
	}

	result.ID = id

	if err := uc.repository.Update(result); err != nil {
		return fmt.Errorf("update ai detect result error: %w", err)
	}

	result, err = uc.repository.GetByIdWithApplication(result.ID)
	if err != nil {
		return fmt.Errorf("get ai detect result by id error(AiDetectStatusUpdateUseCase): %w", err)
	}

	if result.Status == enum.AntiPlagStatusFailed {
		uc.publishReverseEvent(result.JournalID, enum.ServiceAIDetection)
		uc.sendNotification(result.JournalID, result.Application)
		return nil
	}
	if result.Status == enum.AntiPlagStatusSuccess {
		uc.sendNotification(result.JournalID, result.Application)
		return nil
	}

	return &Retry{Message: "RETRY ai detect status update task"}
}

func (uc *AiDetectStatusUpdateUseCase) sendNotification(journalID uint, app *entity.ApplicationEntity) {
	journalMembers, err := uc.roleRepo.GetJournalMemberIds(journalID)
	if err != nil {
		log.Error("get journal members error(AiDetectStatusUpdateUseCase): %w", err)
		return
	}
	if err := uc.notificationService.AiDetectCheckFinished(journalMembers, app.ID, app.Number); err != nil {
		log.Error("send ai detect check finished notification error(AiDetectStatusUpdateUseCase): %w", err)
		return
	}
}

func (uc *AiDetectStatusUpdateUseCase) publishReverseEvent(journalID uint, serviceCode enum.ServiceCode) {
	version := time.Now().UnixMilli()
	eventID := fmt.Sprintf("aidetect_reverse_%d_%d", journalID, version)

	payload, _ := json.Marshal(map[string]any{
		"journal_id":   journalID,
		"service_code": serviceCode.String(),
	})

	event := entity.NewOutboxEventEntity(eventID, version, enum.OutboxEventReverseServiceBalance, payload)
	event, err := uc.outboxRepo.Create(event)
	if err != nil {
		return
	}

	_ = uc.outboxTask.Run(event)
}

type Retry struct {
	Message string
}

func (r *Retry) Error() string {
	return r.Message
}
