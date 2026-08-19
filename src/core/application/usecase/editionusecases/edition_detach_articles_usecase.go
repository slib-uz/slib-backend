package editionusecases

import (
	"context"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type EditionDetachArticlesUseCase struct {
	repository       repository.EditionRepository
	memberPermission *permissionusecases.JournalMemberPermissionUseCase
}

// @inject
func NewEditionDetachArticlesUseCase(
	repository repository.EditionRepository,
	memberPermission *permissionusecases.JournalMemberPermissionUseCase,
) *EditionDetachArticlesUseCase {
	return &EditionDetachArticlesUseCase{
		repository:       repository,
		memberPermission: memberPermission,
	}
}

func (this *EditionDetachArticlesUseCase) Execute(ctx context.Context, editionID uint, articleIDs []uint, user *entity.UserBasicEntity) (int64, error) {
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

	return this.repository.DetachArticles(ctx, editionID, articleIDs)
}
