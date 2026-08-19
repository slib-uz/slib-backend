package aidetectusecases

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/tasks"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
)

type CheckAiDetectUseCase struct {
	gateway          gateway.AiDetectGateway
	aiDetectRepo     repository.AiDetectRepository
	storage          storage.FileStorage
	permission       *permissionusecases.ApplicationReviewerPermissionUseCase
	statusUpdateTask *tasks.AiDetectStatusUpdateTask
	reviewStageRepo  repository.ReviewStageRepository
	billingGateway   gateway.SlibBillingGateway
	outboxRepo       repository.OutboxEventRepository
	outboxTask       *tasks.OutboxEventTaskSender
}

// @inject
func NewCheckAiDetectUseCase(
	gateway gateway.AiDetectGateway,
	aiDetectRepo repository.AiDetectRepository,
	storage storage.FileStorage,
	permission *permissionusecases.ApplicationReviewerPermissionUseCase,
	statusUpdateTask *tasks.AiDetectStatusUpdateTask,
	reviewStageRepo repository.ReviewStageRepository,
	billingGateway gateway.SlibBillingGateway,
	outboxRepo repository.OutboxEventRepository,
	outboxTask *tasks.OutboxEventTaskSender,
) *CheckAiDetectUseCase {
	return &CheckAiDetectUseCase{
		gateway:          gateway,
		aiDetectRepo:     aiDetectRepo,
		storage:          storage,
		permission:       permission,
		statusUpdateTask: statusUpdateTask,
		reviewStageRepo:  reviewStageRepo,
		billingGateway:   billingGateway,
		outboxRepo:       outboxRepo,
		outboxTask:       outboxTask,
	}
}

func (this *CheckAiDetectUseCase) Execute(user *entity.UserBasicEntity, stageID uint) error {
	reviewStage, err := this.reviewStageRepo.GetByIDWithArticle(stageID)
	if err != nil {
		return err
	}

	applicationID := reviewStage.ApplicationID
	journalID := reviewStage.Application.JournalID
	reviewStageID := reviewStage.ID
	articleID := reviewStage.Application.ArticleID

	if err := this.checkAccess(user.Roles, reviewStage.ApplicationID); err != nil {
		return err
	}

	balance, err := this.billingGateway.GetJournalServiceBalance(journalID, enum.ServiceAIDetection)
	if err != nil {
		return err
	}
	if !balance.IsAvailable() {
		return response.InsufficientBalance
	}

	this.publishOutboxEvent(journalID, enum.ServiceAIDetection)

	objectPath := reviewStage.Application.Article.ContentFile

	fileBytes, err := this.storage.GetObject(enum.BucketPrivate, objectPath)
	if err != nil {
		return err
	}

	resultEntity, err := this.gateway.Check(fileBytes, filepath.Base(objectPath))
	if err != nil {
		return err
	}

	resultEntity.ApplicationID = applicationID
	resultEntity.JournalID = journalID
	resultEntity.ReviewStageID = reviewStageID
	resultEntity.ArticleID = articleID

	if err := this.aiDetectRepo.Create(resultEntity); err != nil {
		return err
	}

	return this.statusUpdateTask.Run(resultEntity.ExternalID)
}

func (this *CheckAiDetectUseCase) checkAccess(userRoles []*entity.UserRoleEntity, applicationID uint) error {
	access, err := this.permission.Execute(userRoles, applicationID)
	if err != nil {
		return err
	} else if !access {
		return response.PermissionDeniedError
	}
	return nil
}

func (this *CheckAiDetectUseCase) publishOutboxEvent(journalID uint, serviceCode enum.ServiceCode) {
	version := time.Now().UnixMilli()
	eventID := fmt.Sprintf("aidetect_%d_%d", journalID, version)

	payload, _ := json.Marshal(map[string]any{
		"journal_id":   journalID,
		"service_code": serviceCode.String(),
	})

	event := entity.NewOutboxEventEntity(eventID, version, enum.OutboxEventUseServiceBalance, payload)
	event, err := this.outboxRepo.Create(event)
	if err != nil {
		return
	}

	_ = this.outboxTask.Run(event)
}
