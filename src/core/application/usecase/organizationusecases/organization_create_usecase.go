package organizationusecases

import (
	"errors"

	"slib.uz/src/core/application/response"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type OrganizationCreateUseCase struct {
	repo repository.OrganizationRepository
}

// @inject
func NewOrganizationCreateUseCase(repo repository.OrganizationRepository) *OrganizationCreateUseCase {
	return &OrganizationCreateUseCase{repo: repo}
}

func (this *OrganizationCreateUseCase) Execute(user *entity2.UserBasicEntity, entity *entity2.OrganizationEntity) error {
	if err := this.checkUserAccess(user); err != nil {
		return err
	}

	normalized, err := normalizeAndValidateOrganizationTin(entity.Tin)
	if err != nil {
		return err
	}
	entity.Tin = &normalized

	if existing, err := this.repo.GetByTin(normalized); err == nil {
		return organizationAlreadyExists(existing)
	} else if !errors.Is(err, response.NotFoundError) {
		return err
	}

	if err := this.repo.Create(entity); err != nil {
		if errors.Is(err, response.ConflictError) {
			if existing, lookupErr := this.repo.GetByTin(normalized); lookupErr == nil {
				return organizationAlreadyExists(existing)
			}
		}
		return err
	}
	return nil
}

func (this *OrganizationCreateUseCase) checkUserAccess(user *entity2.UserBasicEntity) error {
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
