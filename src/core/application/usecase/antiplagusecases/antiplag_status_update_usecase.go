package antiplagusecases

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
	"slib.uz/src/core/domain/ports/storage"
)

type AntiPlagStatusUpdateUseCase struct {
	repository          repository.AntiPlagRepository
	gateway             gateway.AntiPlagGateway
	notificationService *service.SendNotificationService
	roleRepo            repository.UserRoleRepository
	storage             storage.FileStorage
	outboxRepo          repository.OutboxEventRepository
	outboxTask          *tasks.OutboxEventTaskSender
}

// @inject
func NewAntiPlagStatusUpdateUseCase(
	repository repository.AntiPlagRepository,
	gateway gateway.AntiPlagGateway,
	notificationService *service.SendNotificationService,
	roleRepo repository.UserRoleRepository,
	storage storage.FileStorage,
	outboxRepo repository.OutboxEventRepository,
	outboxTask *tasks.OutboxEventTaskSender,
) *AntiPlagStatusUpdateUseCase {
	return &AntiPlagStatusUpdateUseCase{
		repository:          repository,
		gateway:             gateway,
		notificationService: notificationService,
		roleRepo:            roleRepo,
		storage:             storage,
		outboxRepo:          outboxRepo,
		outboxTask:          outboxTask,
	}
}

func (this *AntiPlagStatusUpdateUseCase) Execute(externalID uint) error {
	result, err := this.gateway.GetResult(externalID)
	if err != nil {
		return fmt.Errorf("anti plag status update task result error: %w", err)
	}

	id, err := this.repository.FindIDByExternalID(externalID)
	if err != nil {
		return fmt.Errorf("find id from repository error: %w", err)
	}

	result.ID = id

	if err := this.repository.Update(result); err != nil {
		return fmt.Errorf("update anti plag result error: %w", err)
	}

	result, err = this.repository.GetByIdWithApplication(result.ID)
	if err != nil {
		return fmt.Errorf("get anti plag result by id error(AntiPlagStatusUpdateUseCase): %w", err)
	}

	if result.Status == enum.AntiPlagStatusFailed {
		this.publishReverseEvent(result.JournalID, enum.ServiceAntiPlag)
		this.sendNotification(result.JournalID, result.Application)
		return nil
	}
	if result.Status == enum.AntiPlagStatusReady {
		this.sendNotification(result.JournalID, result.Application)
		this.generateCertificate(externalID, result.Application.User.FullName)
		return nil
	}

	return &Retry{Message: "RETRY anti plag status update task"}
}

func (this *AntiPlagStatusUpdateUseCase) generateCertificate(externalID uint, authorFullName string) {
	fileName := fmt.Sprintf("%d.pdf", externalID)

	certificate, err := this.gateway.GenerateCertificate(externalID, authorFullName, fileName)
	if err != nil {
		log.Error("generate certificate error(AntiPlagStatusUpdateUseCase): %w", err)
		return
	}

	objectPath, err := this.storage.PutObject(
		enum.BucketPublic,
		enum.FolderAntiPlagCertificate,
		fileName,
		certificate,
	)
	if err != nil {
		log.Error("put object to storage error(AntiPlagStatusUpdateUseCase): %w", err)
		return
	}

	if err := this.repository.UpdateCertificate(externalID, &objectPath); err != nil {
		log.Error("update certificate in repository error(AntiPlagStatusUpdateUseCase): %w", err)
		return
	}
}

func (this *AntiPlagStatusUpdateUseCase) sendNotification(journalID uint, app *entity.ApplicationEntity) {
	journalMembers, err := this.roleRepo.GetJournalMemberIds(journalID)
	if err != nil {
		log.Error("get journal members error(AntiPlagStatusUpdateUseCase): %w", err)
		return
	}
	if err := this.notificationService.AntiPlagCheckFinished(journalMembers, app.ID, app.Number); err != nil {
		log.Error("send anti plag check finished notification error(AntiPlagStatusUpdateUseCase): %w", err)
		return
	}
}

func (this *AntiPlagStatusUpdateUseCase) publishReverseEvent(journalID uint, serviceCode enum.ServiceCode) {
	version := time.Now().UnixMilli()
	eventID := fmt.Sprintf("antiplag_reverse_%d_%d", journalID, version)

	payload, _ := json.Marshal(map[string]any{
		"journal_id":   journalID,
		"service_code": serviceCode.String(),
	})

	event := entity.NewOutboxEventEntity(eventID, version, enum.OutboxEventReverseServiceBalance, payload)
	event, err := this.outboxRepo.Create(event)
	if err != nil {
		return
	}

	_ = this.outboxTask.Run(event)
}

type Retry struct {
	Message string
}

func (r *Retry) Error() string {
	return r.Message
}
