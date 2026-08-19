package usecase

import (
	"errors"
	"time"

	"github.com/labstack/gommon/log"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/utils"
	"slib.uz/src/infrastructure/config"
)

type PeerReviewStageUseCase struct {
	reviewStageRepo         repository.ReviewStageRepository
	journalMemberPermission *permissionusecases.JournalMemberPermissionUseCase
	journalRepo             repository.JournalRepository
	applicationRepo         repository.ApplicationRepository
	//requirementRepo         repository.ArticleRequirementRepository
	sendNotificationService *service.SendNotificationService
	roleRepo                repository.UserRoleRepository
	env                     *config.Config
}

// @inject
func NewPeerReviewStageUseCase(reviewStageRepo repository.ReviewStageRepository, journalMemberPermission *permissionusecases.JournalMemberPermissionUseCase, journalRepo repository.JournalRepository, applicationRepo repository.ApplicationRepository, sendNotificationService *service.SendNotificationService, roleRepo repository.UserRoleRepository, env *config.Config) *PeerReviewStageUseCase {
	return &PeerReviewStageUseCase{reviewStageRepo: reviewStageRepo, journalMemberPermission: journalMemberPermission, journalRepo: journalRepo, applicationRepo: applicationRepo, sendNotificationService: sendNotificationService, roleRepo: roleRepo, env: env}
}

func (this *PeerReviewStageUseCase) Execute(user *entity.UserBasicEntity, applicationID uint, status enum.Status, reason *string) error {

	reviewStage, err := this.reviewStageRepo.GetBy(applicationID, enum.PeerReviewStage, enum.StatusPending)
	if errors.Is(err, response.NotFoundError) {
		return response.NewFailResponse(400, "application is not in peer review stage")
	}
	if err != nil {
		return err
	}

	// check user permissions this journal
	allowed, err := this.journalMemberPermission.Execute(user.Roles, reviewStage.Application.JournalID)
	if err != nil {
		return err
	}
	if !allowed {
		return response.NewFailResponse(403, "you do not have permission to review this application")
	}

	if err := this.update(user, reviewStage, status, reason); err != nil {
		return err
	}

	this.sendNotification(reviewStage.Application, enum.Status(status), reason)
	return nil
}

func (this *PeerReviewStageUseCase) sendNotification(application *entity.ApplicationEntity, status enum.Status, reason *string) {
	journalMemberIds, err := this.roleRepo.GetJournalMemberIds(application.JournalID)
	if err != nil {
		log.Error("Error sending notification for journal id (PeerReviewStageUseCase):", err)
		return
	}

	receiverIds := append(journalMemberIds, application.UserID)

	senderFunc := this.sendNotificationService.AcceptedAtPeerReview
	if status == enum.StatusRejected {
		senderFunc = this.sendNotificationService.RejectedAtPeerReview
	}

	reasonText := ""
	if reason != nil {
		reasonText = *reason
	}

	if err := senderFunc(receiverIds, application.ID, application.Number, reasonText); err != nil {
		log.Error("Error sending notification (PeerReviewStageUseCase):", err)
	}

}

func (this *PeerReviewStageUseCase) update(user *entity.UserBasicEntity, reviewStage *entity.ReviewStageEntity, status enum.Status, reason *string) error {
	now := time.Now()
	reviewStage.Status = enum.Status(status)
	reviewStage.ReviewerID = &user.ID
	reviewStage.Reason = reason
	reviewStage.ReviewedAt = &now

	// next stage is payment or publish
	application, err := this.applicationRepo.FindByID(reviewStage.ApplicationID)
	if err != nil {
		return err
	}

	journal, err := this.journalRepo.FindByID(application.JournalID)
	if err != nil {
		return err
	}

	nextStage := enum.PublishStage

	// uz: "Agar jurnal ochiq kirish turida bo'lsa maqolani nashr qilish to'lovini muallif to'lashi kerak"
	// en: If the journal is in open access type, the author must pay for the publication
	if journal.AccessType == enum.PublicAccessType && journal.PublishingPrice != nil && *journal.PublishingPrice > 0 {
		nextStage = enum.PaymentStage
	}

	var nextReviewStage *entity.ReviewStageEntity

	switch status {
	case enum.StatusAccepted:
		nextReviewStage = entity.NewReviewStageEntity(
			0,
			reviewStage.ApplicationID,
			nil,
			nextStage,
			enum.StatusPending,
			nil,
			nil,
			nil,
			nil,
			time.Now(),
			nil,
			nil,
			false,
			nil,
			nil,
			nil,
			utils.MakeDeadline(time.Now(), this.env.ReviewDeadlineDays),
			nil,
		)
	case enum.StatusRejected:
		resubmitDeadline := utils.MakeDeadline(time.Now(), this.env.ResubmitDeadlineDays)
		reviewStage.ResubmitDeadline = &resubmitDeadline
	}

	return this.reviewStageRepo.UpdateToNext(reviewStage, nextReviewStage)
}
