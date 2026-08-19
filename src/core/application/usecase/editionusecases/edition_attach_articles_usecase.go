package editionusecases

import (
	"context"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type EditionAttachArticlesUseCase struct {
	repository       repository.EditionRepository
	memberPermission *permissionusecases.JournalMemberPermissionUseCase
}

// @inject
func NewEditionAttachArticlesUseCase(repository repository.EditionRepository, memberPermission *permissionusecases.JournalMemberPermissionUseCase) *EditionAttachArticlesUseCase {
	return &EditionAttachArticlesUseCase{repository: repository, memberPermission: memberPermission}
}

func (this *EditionAttachArticlesUseCase) Execute(ctx context.Context, editionID uint, articleIDs []uint, user *entity.UserBasicEntity) (int64, error) {
	edition, err := this.repository.GetByID(ctx, editionID)
	if err != nil {
		return 0, err
	}

	allowed, err := this.memberPermission.Execute(user.Roles, edition.JournalID)
	if err != nil {
		return 0, err
	}
	if !allowed {
		return 0, response.PermissionDeniedError
	}

	return this.repository.AttachArticles(ctx, editionID, articleIDs)
}
