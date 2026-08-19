package organizationusecases

import (
	"errors"

	"slib.uz/src/core/application/response"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type OrganizationUpdateUseCase struct {
	repo repository.OrganizationRepository
}

// @inject
func NewOrganizationUpdateUseCase(repo repository.OrganizationRepository) *OrganizationUpdateUseCase {
	return &OrganizationUpdateUseCase{repo: repo}
}

func (this *OrganizationUpdateUseCase) Execute(id uint, entity *entity2.OrganizationEntity, user *entity2.UserBasicEntity) error {
	if err := this.checkUserAccess(user); err != nil {
		return err
	}

	normalized, err := normalizeAndValidateOrganizationTin(entity.Tin)
	if err != nil {
		return err
	}
	entity.Tin = &normalized

	if existing, err := this.repo.GetByTin(normalized); err == nil {
		if existing.ID != id {
			return organizationAlreadyExists(existing)
		}
	} else if !errors.Is(err, response.NotFoundError) {
		return err
	}

	if err := this.repo.Update(id, entity); err != nil {
		if errors.Is(err, response.ConflictError) {
			if existing, lookupErr := this.repo.GetByTin(normalized); lookupErr == nil && existing.ID != id {
				return organizationAlreadyExists(existing)
			}
		}
		return err
	}
	return nil
}

func (this *OrganizationUpdateUseCase) checkUserAccess(user *entity2.UserBasicEntity) error {
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
