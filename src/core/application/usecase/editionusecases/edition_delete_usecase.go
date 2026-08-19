package editionusecases

import (
	"context"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/permissionusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type EditionDeleteUseCase struct {
	repository       repository.EditionRepository
	articleRepo      repository.ArticleRepository
	memberPermission *permissionusecases.JournalMemberPermissionUseCase
}

// @inject
func NewEditionDeleteUseCase(
	repository repository.EditionRepository,
	articleRepo repository.ArticleRepository,
	memberPermission *permissionusecases.JournalMemberPermissionUseCase,
) *EditionDeleteUseCase {
	return &EditionDeleteUseCase{
		repository:       repository,
		articleRepo:      articleRepo,
		memberPermission: memberPermission,
	}
}

func (this *EditionDeleteUseCase) Execute(ctx context.Context, editionID uint, user *entity.UserBasicEntity) error {
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

	articleCount, err := this.articleRepo.GetArticleCountByEditionID(editionID)
	if err != nil {
		return err
	}

	if articleCount > 0 {
		return response.EditionHasArticles
	}

	return this.repository.Delete(ctx, editionID)
}
