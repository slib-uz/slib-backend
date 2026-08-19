package aidetectusecases

import (
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type AiDetectResultListUseCase struct {
	repository                    repository.AiDetectRepository
	applicationReviewerPermission *permissionusecases.ApplicationReviewerPermissionUseCase
	applicationOwnerPermission    *permissionusecases.ApplicationOwnerPermissionUseCase
}

// @inject
func NewAiDetectResultListUseCase(repository repository.AiDetectRepository, applicationReviewerPermission *permissionusecases.ApplicationReviewerPermissionUseCase, applicationOwnerPermission *permissionusecases.ApplicationOwnerPermissionUseCase) *AiDetectResultListUseCase {
	return &AiDetectResultListUseCase{repository: repository, applicationReviewerPermission: applicationReviewerPermission, applicationOwnerPermission: applicationOwnerPermission}
}

func (this *AiDetectResultListUseCase) Execute(user *entity.UserBasicEntity, applicationID uint, reviewStageID uint) ([]*entity.AiDetectResultEntity, error) {

	if err := this.checkAccess(user, applicationID); err != nil {
		return nil, err
	}

	results, err := this.repository.GetAllByReviewStageID(reviewStageID)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (this *AiDetectResultListUseCase) checkAccess(user *entity.UserBasicEntity, applicationID uint) error {
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
