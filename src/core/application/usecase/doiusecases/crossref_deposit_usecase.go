package doiusecases

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/tasks"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/repository"
)

type CrossRefDepositUseCase struct {
	crossRefGateway gateway.CrossRefGateway
	billingGateway  gateway.SlibBillingGateway
	doiDepositRepo  repository.DoiDepositRepository
	doiSettingRepo  repository.JournalDoiSettingRepository
	articleRepo     repository.ArticleRepository
	journalRepo     repository.JournalRepository
	outboxRepo      repository.OutboxEventRepository
	outboxTask      *tasks.OutboxEventTaskSender
}

// @inject
func NewCrossRefDepositUseCase(
	crossRefGateway gateway.CrossRefGateway,
	billingGateway gateway.SlibBillingGateway,
	doiDepositRepo repository.DoiDepositRepository,
	doiSettingRepo repository.JournalDoiSettingRepository,
	articleRepo repository.ArticleRepository,
	journalRepo repository.JournalRepository,
	outboxRepo repository.OutboxEventRepository,
	outboxTask *tasks.OutboxEventTaskSender,
) *CrossRefDepositUseCase {
	return &CrossRefDepositUseCase{
		crossRefGateway: crossRefGateway,
		billingGateway:  billingGateway,
		doiDepositRepo:  doiDepositRepo,
		doiSettingRepo:  doiSettingRepo,
		articleRepo:     articleRepo,
		journalRepo:     journalRepo,
		outboxRepo:      outboxRepo,
		outboxTask:      outboxTask,
	}
}

func (this *CrossRefDepositUseCase) Execute(ctx context.Context, articleID uint) error {
	article, err := this.articleRepo.FindByIDWithAuthors(articleID)
	if err != nil {
		return err
	}

	journal, err := this.journalRepo.GetByIDWithRelations(article.JournalID)
	if err != nil {
		return err
	}

	// Check service balance before deposit
	balance, err := this.billingGateway.GetJournalServiceBalance(journal.ID, enum.ServiceDOI)
	if err != nil {
		return err
	}
	if !balance.IsAvailable() {
		return response.InsufficientBalance
	}

	setting, err := this.doiSettingRepo.GetByJournalID(ctx, journal.ID)
	if err != nil {
		return response.NewOptionalResponse(200, response.CodeNotReady, nil, fmt.Sprintf("journal %d has no DOI settings configured", journal.ID))
	}

	if article.DOI == nil || *article.DOI == "" {
		doi := setting.GenerateDOI(article.ID)
		article.DOI = &doi
	}

	batchID := fmt.Sprintf("batch_%d_%d", journal.ID, article.ID)

	result, depositErr := this.crossRefGateway.Deposit(&gateway.CrossRefDepositParams{
		BatchID:     batchID,
		Username:    setting.Username,
		Password:    setting.Password,
		JournalName: setting.JournalName,
		Journal:     journal,
		Articles:    []*entity.ArticleEntity{article},
	})

	this.saveDeposit(ctx, articleID, batchID, *article.DOI, result)

	if result != nil && result.Status == enum.DoiDepositStatusSuccess && result.DOI != "" {
		_ = this.articleRepo.UpdateDOI(articleID, result.DOI)
		// Deduct service balance via outbox pattern
		this.publishOutboxEvent(journal.ID, enum.ServiceDOI)
	}

	if depositErr != nil {
		return depositErr
	}
	if result != nil && result.Status == enum.DoiDepositStatusFailure {
		return response.NewFailResponse(400, result.Message)
	}

	return nil
}

func (this *CrossRefDepositUseCase) publishOutboxEvent(journalID uint, serviceCode enum.ServiceCode) {
	eventType := enum.OutboxEventUseServiceBalance
	version := time.Now().UnixMilli()
	eventID := fmt.Sprintf("doi_%d_%d", journalID, version)

	payload, _ := json.Marshal(map[string]any{
		"journal_id":   journalID,
		"service_code": serviceCode.String(),
	})

	event := entity.NewOutboxEventEntity(eventID, version, eventType, payload)
	event, err := this.outboxRepo.Create(event)
	if err != nil {
		return
	}

	_ = this.outboxTask.Run(event)
}

func (this *CrossRefDepositUseCase) saveDeposit(ctx context.Context, articleID uint, batchID string, doi string, result *entity.CrossRefDepositResultEntity) {
	status, message, submissionID, requestBody, responseBody := enum.DoiDepositStatusFailure, "", "", "", ""
	if result != nil {
		status = result.Status
		message = result.Message
		submissionID = result.SubmissionID
		requestBody = result.RequestBody
		responseBody = result.ResponseBody
		if result.DOI != "" {
			doi = result.DOI
		}
	}

	existing, _ := this.doiDepositRepo.GetByBatchID(ctx, batchID)
	if existing != nil {
		_ = this.doiDepositRepo.UpdateByBatchID(ctx, batchID, status, message, submissionID, requestBody, responseBody)
	} else {
		deposit := entity.NewDoiDepositEntity(0, &articleID, batchID, doi, status, message, submissionID, requestBody, responseBody, time.Time{}, time.Time{})
		_ = this.doiDepositRepo.Create(ctx, deposit)
	}
}
