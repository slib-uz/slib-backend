package authorshipclaimusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type ListAuthorshipClaimsUseCase struct {
	repo repository.AuthorshipClaimRepository
}

// @inject
func NewListAuthorshipClaimsUseCase(repo repository.AuthorshipClaimRepository) *ListAuthorshipClaimsUseCase {
	return &ListAuthorshipClaimsUseCase{repo: repo}
}

func (this *ListAuthorshipClaimsUseCase) Execute(page, size int, filters map[string]interface{}) (*entity.PagingEntity[entity.AuthorshipClaimEntity], error) {
	return this.repo.GetList(page, size, filters)
}
