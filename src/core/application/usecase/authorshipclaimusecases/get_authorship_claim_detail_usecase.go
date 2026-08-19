package authorshipclaimusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type GetAuthorshipClaimDetailUseCase struct {
	repo repository.AuthorshipClaimRepository
}

// @inject
func NewGetAuthorshipClaimDetailUseCase(repo repository.AuthorshipClaimRepository) *GetAuthorshipClaimDetailUseCase {
	return &GetAuthorshipClaimDetailUseCase{repo: repo}
}

func (this *GetAuthorshipClaimDetailUseCase) Execute(id uint) (*entity.AuthorshipClaimEntity, error) {
	return this.repo.GetByID(id)
}
