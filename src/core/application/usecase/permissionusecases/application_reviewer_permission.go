package permissionusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ApplicationReviewerPermissionUseCase struct {
	journalMemberPermission *JournalMemberPermissionUseCase
	applicationRepository   repository.ApplicationRepository
}

// @inject
func NewApplicationReviewerPermissionUseCase(journalMemberPermission *JournalMemberPermissionUseCase, applicationRepository repository.ApplicationRepository) *ApplicationReviewerPermissionUseCase {
	return &ApplicationReviewerPermissionUseCase{journalMemberPermission: journalMemberPermission, applicationRepository: applicationRepository}
}

func (this *ApplicationReviewerPermissionUseCase) Execute(userRoles []*entity.UserRoleEntity, applicationID uint) (bool, error) {
	application, err := this.applicationRepository.FindByID(applicationID)
	if err != nil {
		return false, err
	}
	return this.journalMemberPermission.Execute(userRoles, application.JournalID)
}
