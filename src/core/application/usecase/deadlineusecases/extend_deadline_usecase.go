package deadlineusecases

import (
	"time"

	"slib.uz/src/core/application/response"
	permissionusecases2 "slib.uz/src/core/application/usecase/permissionusecases"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type ExtendDeadlineUseCase struct {
	reviewStageRepo         repository.ReviewStageRepository
	roleRepo                repository.UserRoleRepository
	journalMemberPermission *permissionusecases2.JournalMemberPermissionUseCase
	chiefEditorPermission   *permissionusecases2.ChiefEditorPermissionUseCase
}

// @inject
func NewExtendDeadlineUseCase(reviewStageRepo repository.ReviewStageRepository, roleRepo repository.UserRoleRepository) *ExtendDeadlineUseCase {
	return &ExtendDeadlineUseCase{reviewStageRepo: reviewStageRepo, roleRepo: roleRepo}
}

func (this *ExtendDeadlineUseCase) Execute(user *entity2.UserBasicEntity, reviewStageID uint, newDeadline time.Time, deadlineType enum.DeadlineType) error {

	reviewStage, err := this.reviewStageRepo.GetWithApplication(reviewStageID)
	if err != nil {
		return err
	}

	if err := this.hasPermission(user, reviewStage.Application.JournalID, deadlineType); err != nil {
		return err
	}

	if err := this.checkDeadline(reviewStage, newDeadline, deadlineType); err != nil {
		return err
	}
	return this.reviewStageRepo.UpdateDeadline(reviewStageID, newDeadline, enum.DeadlineType(deadlineType))
}

func (this *ExtendDeadlineUseCase) checkDeadline(reviewStage *entity2.ReviewStageEntity, newDeadline time.Time, deadlineType enum.DeadlineType) error {
	now := time.Now()

	oldDeadline := reviewStage.Deadline
	if deadlineType == enum.DeadlineTypeResubmitDeadline {
		oldDeadline = *reviewStage.ResubmitDeadline
	}

	if now.Before(oldDeadline) {
		return response.NewFailResponse(400, "Deadline has not yet passed")
	}

	if newDeadline.Before(oldDeadline) {
		return response.NewFailResponse(400, "New deadline must be after the old deadline")
	}

	return nil
}

func (this *ExtendDeadlineUseCase) hasPermission(user *entity2.UserBasicEntity, journalID uint, deadlineType enum.DeadlineType) error {

	if deadlineType == enum.DeadlineTypeReviewDeadline && this.chiefEditorPermission.Execute(user.Roles, journalID) {
		return nil
	}

	// Qo'shma shart ikkiga bo'lindi: aks holda && qisqa tutashuvi bilan
	// xatoni qayta ishlash chalkash bo'lardi va deadlineType mos kelmaganda
	// ham DB so'rovi ketishi mumkin edi.
	if deadlineType == enum.DeadlineTypeResubmitDeadline {
		isMember, err := this.journalMemberPermission.Execute(user.Roles, journalID)
		if err != nil {
			return err
		}
		if isMember {
			return nil
		}
	}

	return response.PermissionDeniedError
}

func (this *ExtendDeadlineUseCase) sendNotification() {

}
