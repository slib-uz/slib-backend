package spellcheckusecases

import (
	"slib.uz/src/core/application/mapper"
	"slib.uz/src/core/application/response"
	permissionusecases2 "slib.uz/src/core/application/usecase/permissionusecases"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type SpellCheckResultsListUsecase struct {
	repository                    repository.SpellCheckResultRepository
	applicationReviewerPermission *permissionusecases2.ApplicationReviewerPermissionUseCase
	applicationOwnerPermission    *permissionusecases2.ApplicationOwnerPermissionUseCase
}

// @inject
func NewSpellCheckResultsListUsecase(repository repository.SpellCheckResultRepository, applicationReviewerPermission *permissionusecases2.ApplicationReviewerPermissionUseCase, applicationOwnerPermission *permissionusecases2.ApplicationOwnerPermissionUseCase) *SpellCheckResultsListUsecase {
	return &SpellCheckResultsListUsecase{repository: repository, applicationReviewerPermission: applicationReviewerPermission, applicationOwnerPermission: applicationOwnerPermission}
}

func (this *SpellCheckResultsListUsecase) Execute(user *entity2.UserBasicEntity, applicationID uint, reviewStageID uint) ([]*entity2.SpellCheckResultEntity, error) {

	if err := this.checkAccess(user, applicationID); err != nil {
		return nil, err
	}

	var results []*entity2.SpellCheckResultEntity
	var err error

	if reviewStageID == 0 {
		results, err = this.repository.GetByApplicationID(applicationID)
	} else {
		results, err = this.repository.GetByReviewStageID(reviewStageID)
	}

	if err != nil {
		return nil, err
	}

	var responseEntities = make([]*entity2.SpellCheckResultEntity, len(results))
	for i, item := range results {
		responseEntities[i] = mapper.SpellCheckResultEntityToResponseEntity(item)
	}
	return responseEntities, nil
}

func (this *SpellCheckResultsListUsecase) checkAccess(user *entity2.UserBasicEntity, applicationID uint) error {
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
