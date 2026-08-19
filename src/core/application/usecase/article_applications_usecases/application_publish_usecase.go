package usecase

import (
	"time"

	"github.com/labstack/gommon/log"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/application/tasks"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type ApplicationPublishUseCase struct {
	applicationRepo         repository.ApplicationRepository
	reviewStageRepo         repository.ReviewStageRepository
	permission              *permissionusecases.JournalMemberPermissionUseCase
	sendNotificationService *service.SendNotificationService
	roleRepo                repository.UserRoleRepository
	roiSenderTask           *tasks.SetRoiSenderTask
}

// @inject
func NewApplicationPublishUseCase(
	applicationRepo repository.ApplicationRepository,
	reviewStageRepo repository.ReviewStageRepository,
	permission *permissionusecases.JournalMemberPermissionUseCase,
	sendNotificationService *service.SendNotificationService,
	roleRepo repository.UserRoleRepository,
	roiSenderTask *tasks.SetRoiSenderTask,
) *ApplicationPublishUseCase {
	return &ApplicationPublishUseCase{
		applicationRepo:         applicationRepo,
		reviewStageRepo:         reviewStageRepo,
		permission:              permission,
		sendNotificationService: sendNotificationService,
		roleRepo:                roleRepo,
		roiSenderTask:           roiSenderTask,
	}
}

func (this *ApplicationPublishUseCase) Execute(user *entity.UserBasicEntity, applicationID uint, file string) error {
	application, err := this.applicationRepo.GetWithArticleAndJournal(applicationID)
	if err != nil {
		return err
	}

	// check permission
	allowed, err := this.permission.Execute(user.Roles, application.Article.JournalID)
	if err != nil {
		return err
	}
	if !allowed {
		return response.PermissionDeniedError
	}

	reviewStage, err := this.reviewStageRepo.GetBy(applicationID, enum.PublishStage, enum.StatusPending)
	if err != nil {
		return response.NewFailResponse(400, "Application is not in publish stage in pending status")
	}

	// Update review stage with reviewer info
	now := time.Now()
	reviewStage.ReviewerID = &user.ID
	reviewStage.Status = enum.StatusAccepted
	reviewStage.ReviewedAt = &now

	if err := this.reviewStageRepo.UpdateToNext(reviewStage, nil); err != nil {
		return err
	}

	if err := this.applicationRepo.Publish(application.Article.ID, applicationID, file); err != nil {
		return err
	}

	// send to roi task sender
	if err := this.roiSenderTask.Run(tasks.SetRoiSenderPayload{
		User:          user,
		ApplicationID: applicationID,
		File:          file,
	}); err != nil {
		return err
	}

	this.sendNotification(application)

	return nil
}

func (this *ApplicationPublishUseCase) sendNotification(application *entity.ApplicationEntity) {
	journalMemberIds, err := this.roleRepo.GetJournalMemberIds(application.JournalID)
	if err != nil {
		log.Error("Error sending notification for journal id (ApplicationPublishUseCase):", err)
		return
	}

	receiverIds := append(journalMemberIds, application.UserID)

	if err := this.sendNotificationService.
		ApplicationPublished(receiverIds, application.ID, application.Number); err != nil {
		log.Error("Error sending notification for application publish (ApplicationPublishUseCase):", err)
	}
}
