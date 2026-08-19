package antiplagusecases

import (
	"slib.uz/src/core/application/response"
	permissionusecases2 "slib.uz/src/core/application/usecase/permissionusecases"
	entity2 "slib.uz/src/core/domain/entity"

	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type AntiPlagResultListUseCase struct {
	repository                    repository.AntiPlagRepository
	applicationReviewerPermission *permissionusecases2.ApplicationReviewerPermissionUseCase
	applicationOwnerPermission    *permissionusecases2.ApplicationOwnerPermissionUseCase
}

// @inject
func NewAntiPlagResultListUseCase(repository repository.AntiPlagRepository, applicationReviewerPermission *permissionusecases2.ApplicationReviewerPermissionUseCase, applicationOwnerPermission *permissionusecases2.ApplicationOwnerPermissionUseCase) *AntiPlagResultListUseCase {
	return &AntiPlagResultListUseCase{repository: repository, applicationReviewerPermission: applicationReviewerPermission, applicationOwnerPermission: applicationOwnerPermission}
}

func (this *AntiPlagResultListUseCase) Execute(user *entity2.UserBasicEntity, applicationID uint, reviewStageID uint) ([]*entity2.AntiPlagResultEntity, error) {

	if err := this.checkAccess(user, applicationID); err != nil {
		return nil, err
	}
	var results []*entity2.AntiPlagResultEntity
	var err error

	if reviewStageID == 0 {
		results, err = this.repository.GetAllByApplicationID(applicationID)
	} else {
		results, err = this.repository.GetAllByReviewStageID(reviewStageID)
	}

	if err != nil {
		return nil, err
	}

	return results, nil

}

func (this *AntiPlagResultListUseCase) checkAccess(user *entity2.UserBasicEntity, applicationID uint) error {
	for _, v := range user.Roles {
		if v.Role == enum.RoleAdmin {
			return nil
		}
	}
	if access, err := this.applicationReviewerPermission.Execute(user.Roles, applicationID); err != nil {
		return err
	} else if access {
		return nil
	}
	if access, err := this.applicationOwnerPermission.Execute(user.ID, applicationID); err != nil {
		return err
	} else if access {
		return nil
	}
	return response.PermissionDeniedError
}
