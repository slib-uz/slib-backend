package editionusecases

import (
	"context"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type EditionUpdateUseCase struct {
	repository       repository.EditionRepository
	memberPermission *permissionusecases.JournalMemberPermissionUseCase
}

// @inject
func NewEditionUpdateUseCase(repository repository.EditionRepository, memberPermission *permissionusecases.JournalMemberPermissionUseCase) *EditionUpdateUseCase {
	return &EditionUpdateUseCase{repository: repository, memberPermission: memberPermission}
}

func (this *EditionUpdateUseCase) Execute(ctx context.Context, editionID uint, data *entity.EditionEntity, user *entity.UserBasicEntity) error {
	edition, err := this.repository.GetByID(ctx, editionID)
	if err != nil {
		return err
	}

	allowed, err := this.memberPermission.Execute(user.Roles, edition.JournalID)
	if err != nil {
		return err
	}
	if !allowed {
		return response.PermissionDeniedError
	}

	data.ID = editionID
	return this.repository.Update(ctx, data)
}
