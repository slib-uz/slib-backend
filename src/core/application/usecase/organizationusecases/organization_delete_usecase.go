package organizationusecases

import (
	"slib.uz/src/core/application/response"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type OrganizationDeleteUseCase struct {
	repo repository.OrganizationRepository
}

// @inject
func NewOrganizationDeleteUseCase(repo repository.OrganizationRepository) *OrganizationDeleteUseCase {
	return &OrganizationDeleteUseCase{repo: repo}
}

func (this *OrganizationDeleteUseCase) Execute(user *entity2.UserBasicEntity, id uint) error {
	if err := this.checkUserAccess(user); err != nil {
		return err
	}
	return this.repo.Delete(id)
}

func (this *OrganizationDeleteUseCase) checkUserAccess(user *entity2.UserBasicEntity) error {
	if user.IsAdmin {
		return nil
	}
	for _, role := range user.Roles {
		if role.Role == enum.RoleAdmin {
			return nil
		}
	}
	return response.PermissionDeniedError
}
