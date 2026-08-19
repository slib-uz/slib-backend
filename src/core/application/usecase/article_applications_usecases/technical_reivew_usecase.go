package usecase

import (
	"errors"
	"time"

	"github.com/labstack/gommon/log"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/application/usecase/permissionusecases"
	entity2 "slib.uz/src/core/domain/entity"
	enum2 "slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/utils"
	"slib.uz/src/infrastructure/config"
)

type TechnicalReviewUsecase struct {
	reviewStageRepository   repository.ReviewStageRepository
	journalMemberPermission *permissionusecases.JournalMemberPermissionUseCase
	sendNotificationService *service.SendNotificationService
	roleRepo                repository.UserRoleRepository
	env                     *config.Config
}

// @inject
func NewTechnicalReviewUsecase(reviewStageRepository repository.ReviewStageRepository, journalMemberPermission *permissionusecases.JournalMemberPermissionUseCase, sendNotificationService *service.SendNotificationService, roleRepo repository.UserRoleRepository, env *config.Config) *TechnicalReviewUsecase {
	return &TechnicalReviewUsecase{reviewStageRepository: reviewStageRepository, journalMemberPermission: journalMemberPermission, sendNotificationService: sendNotificationService, roleRepo: roleRepo, env: env}
}

func (this *TechnicalReviewUsecase) Execute(user *entity2.UserBasicEntity, applicationID uint, status enum2.Status, reason *string) error {

	reviewStage, err := this.reviewStageRepository.GetBy(applicationID, enum2.TechnicalReviewStage, enum2.StatusPending)

	if errors.Is(err, response.NotFoundError) {
		return response.NewFailResponse(400, "application is not in technical review stage")
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
		log.Error("Error updating review stage in technical review use case:", err)
	}

	this.sendNotification(reviewStage.Application, enum2.Status(status), reason)
	return nil
}

func (this *TechnicalReviewUsecase) sendNotification(application *entity2.ApplicationEntity, status enum2.Status, reason *string) {
	journalMemberIds, err := this.roleRepo.GetJournalMemberIds(application.JournalID)
	if err != nil {
		log.Error("Error sending notification for journal id (TechnicalReviewUsecase):", err)
		return
	}

	receiverIds := append(journalMemberIds, application.UserID)

	senderFunc := this.sendNotificationService.AcceptedAtTechnicalReview
	if status == enum2.StatusRejected {
		senderFunc = this.sendNotificationService.RejectedAtTechnicalReview
	}

	reasonText := ""
	if reason != nil {
		reasonText = *reason
	}

	if err := senderFunc(receiverIds, application.ID, application.Number, reasonText); err != nil {
		log.Error("Error sending notification (TechnicalReviewUsecase):", err)
	}
}

func (this *TechnicalReviewUsecase) update(user *entity2.UserBasicEntity, reviewStage *entity2.ReviewStageEntity, status enum2.Status, reason *string) error {
	now := time.Now()
	reviewStage.Status = enum2.Status(status)
	reviewStage.ReviewerID = &user.ID
	reviewStage.Reason = reason
	reviewStage.ReviewedAt = &now

	var nextStage *entity2.ReviewStageEntity

	switch status {
	case enum2.StatusAccepted:
		nextStage = entity2.NewReviewStageEntity(
			0,
			reviewStage.ApplicationID,
			nil,
			enum2.PeerReviewStage,
			enum2.StatusPending,
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
	case enum2.StatusRejected:
		resubmitDeadline := utils.MakeDeadline(time.Now(), this.env.ResubmitDeadlineDays)
		reviewStage.ResubmitDeadline = &resubmitDeadline
	}

	return this.reviewStageRepository.UpdateToNext(reviewStage, nextStage)
}
