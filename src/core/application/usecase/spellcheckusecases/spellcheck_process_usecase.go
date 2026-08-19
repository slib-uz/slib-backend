package spellcheckusecases

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/application/tasks"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
	"slib.uz/src/core/utils"
)

type SpellCheckProcessUseCase struct {
	literacyGateway     gateway.LiteracyGateway
	storage             storage.FileStorage
	resultRepository    repository.SpellCheckResultRepository
	notificationService *service.SendNotificationService
	roleRepo            repository.UserRoleRepository
	billingGateway      gateway.SlibBillingGateway
	outboxRepo          repository.OutboxEventRepository
	outboxTask          *tasks.OutboxEventTaskSender
}

// @inject
func NewSpellCheckProcessUseCase(
	literacyGateway gateway.LiteracyGateway,
	storage storage.FileStorage,
	resultRepository repository.SpellCheckResultRepository,
	notificationService *service.SendNotificationService,
	roleRepo repository.UserRoleRepository,
	billingGateway gateway.SlibBillingGateway,
	outboxRepo repository.OutboxEventRepository,
	outboxTask *tasks.OutboxEventTaskSender,
) *SpellCheckProcessUseCase {
	return &SpellCheckProcessUseCase{
		literacyGateway:     literacyGateway,
		storage:             storage,
		resultRepository:    resultRepository,
		notificationService: notificationService,
		roleRepo:            roleRepo,
		billingGateway:      billingGateway,
		outboxRepo:          outboxRepo,
		outboxTask:          outboxTask,
	}
}

func (uc *SpellCheckProcessUseCase) Execute(payload *entity.SpellCheckResultEntity) error {
	file, err := uc.storage.GetObject(enum.BucketPrivate, payload.File)
	if err != nil {
		log.Printf("Error getting file for spellcheck: %v\n", err)
		_ = uc.setResult(payload.ID, -1, nil, nil)
		return err
	}

	resultFile, err := uc.literacyGateway.SpellCheck(file)
	if err != nil {
		log.Printf("Error during spellcheck: %v\n", err)
		_ = uc.setResult(payload.ID, -1, nil, nil)
		return err
	}

	balance, err := uc.billingGateway.GetJournalServiceBalance(payload.JournalID, enum.ServiceSavodxon)
	if err != nil {
		log.Printf("Error checking journal service balance: %v\n", err)
		_ = uc.setResult(payload.ID, -1, nil, nil)
		return err
	}
	if !balance.IsAvailable() {
		_ = uc.setResult(payload.ID, -1, nil, nil)
		return response.InsufficientBalance
	}

	uc.publishOutboxEvent(payload.JournalID, enum.ServiceSavodxon)

	path, err := uc.storage.PutObject(
		enum.BucketPublic,
		enum.FolderSpellCheck,
		utils.MakeFileName(payload.File, utils.RandomHex(1)),
		resultFile,
	)
	if err != nil {
		log.Printf("Error saving spellcheck result: %v\n", err)
		_ = uc.setResult(payload.ID, -1, nil, nil)
		uc.sendNotification(payload.ID)
		return err
	}

	uc.sendNotification(payload.ID)

	now := time.Now()
	return uc.setResult(payload.ID, 1, &path, &now)
}

func (uc *SpellCheckProcessUseCase) setResult(id uint, status int, resultFile *string, resultTime *time.Time) error {
	return uc.resultRepository.Update(id, status, resultFile, resultTime)
}

func (uc *SpellCheckProcessUseCase) sendNotification(id uint) {
	result, err := uc.resultRepository.GetByIDWithApplication(id)
	if err != nil {
		log.Printf("Error getting spellcheck result by id: %v\n", err)
		return
	}
	app := result.Application

	journalMembers, err := uc.roleRepo.GetJournalMemberIds(result.JournalID)
	if err != nil {
		log.Printf("Error getting journal members: %v\n", err)
		return
	}
	if err := uc.notificationService.SpellCheckFinished(journalMembers, app.ID, app.Number); err != nil {
		log.Printf("Error sending spellcheck finished notification: %v\n", err)
		return
	}
}

func (uc *SpellCheckProcessUseCase) publishOutboxEvent(journalID uint, serviceCode enum.ServiceCode) {
	version := time.Now().UnixMilli()
	eventID := fmt.Sprintf("spellcheck_%d_%d", journalID, version)

	payload, _ := json.Marshal(map[string]any{
		"journal_id":   journalID,
		"service_code": serviceCode.String(),
	})

	event := entity.NewOutboxEventEntity(eventID, version, enum.OutboxEventUseServiceBalance, payload)
	event, err := uc.outboxRepo.Create(event)
	if err != nil {
		return
	}

	_ = uc.outboxTask.Run(event)
}
