package antiplagusecases

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/tasks"
	"slib.uz/src/core/application/usecase/permissionusecases"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/storage"
)

type CheckAntiPlagUseCase struct {
	gateway            gateway.AntiPlagGateway
	antiPlagResultRepo repository.AntiPlagRepository
	applicationRepo    repository.ApplicationRepository
	storage            storage.FileStorage
	permission         *permissionusecases.ApplicationReviewerPermissionUseCase
	statusUpdateTask   *tasks.AntiPlagStatusUpdateTask
	reviewStageRepo    repository.ReviewStageRepository
	billingGateway     gateway.SlibBillingGateway
	outboxRepo         repository.OutboxEventRepository
	outboxTask         *tasks.OutboxEventTaskSender
}

// @inject
func NewCheckAntiPlagUseCase(
	gateway gateway.AntiPlagGateway,
	antiPlagResultRepo repository.AntiPlagRepository,
	applicationRepo repository.ApplicationRepository,
	storage storage.FileStorage,
	permission *permissionusecases.ApplicationReviewerPermissionUseCase,
	statusUpdateTask *tasks.AntiPlagStatusUpdateTask,
	reviewStageRepo repository.ReviewStageRepository,
	billingGateway gateway.SlibBillingGateway,
	outboxRepo repository.OutboxEventRepository,
	outboxTask *tasks.OutboxEventTaskSender,
) *CheckAntiPlagUseCase {
	return &CheckAntiPlagUseCase{
		gateway:            gateway,
		antiPlagResultRepo: antiPlagResultRepo,
		applicationRepo:    applicationRepo,
		storage:            storage,
		permission:         permission,
		statusUpdateTask:   statusUpdateTask,
		reviewStageRepo:    reviewStageRepo,
		billingGateway:     billingGateway,
		outboxRepo:         outboxRepo,
		outboxTask:         outboxTask,
	}
}

func (this *CheckAntiPlagUseCase) Execute(user *entity2.UserBasicEntity, stageID uint) error {

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

	balance, err := this.billingGateway.GetJournalServiceBalance(journalID, enum.ServiceAntiPlag)
	if err != nil {
		return err
	}
	if !balance.IsAvailable() {
		return response.InsufficientBalance
	}

	this.publishOutboxEvent(journalID, enum.ServiceAntiPlag)

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

	if err := this.antiPlagResultRepo.Create(resultEntity); err != nil {
		return err
	}

	return this.statusUpdateTask.Run(resultEntity.ExternalID)
}

func (this *CheckAntiPlagUseCase) checkAccess(userRoles []*entity2.UserRoleEntity, applicationID uint) error {
	access, err := this.permission.Execute(userRoles, applicationID)
	if err != nil {
		return err
	} else if !access {
		return response.PermissionDeniedError
	}
	return nil
}

func (this *CheckAntiPlagUseCase) publishOutboxEvent(journalID uint, serviceCode enum.ServiceCode) {
	version := time.Now().UnixMilli()
	eventID := fmt.Sprintf("antiplag_%d_%d", journalID, version)

	payload, _ := json.Marshal(map[string]any{
		"journal_id":   journalID,
		"service_code": serviceCode.String(),
	})

	event := entity2.NewOutboxEventEntity(eventID, version, enum.OutboxEventUseServiceBalance, payload)
	event, err := this.outboxRepo.Create(event)
	if err != nil {
		return
	}

	_ = this.outboxTask.Run(event)
}
