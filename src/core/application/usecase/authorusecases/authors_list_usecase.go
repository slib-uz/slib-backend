package authorusecases

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type AuthorsListUsecase struct {
	repository repository.AuthorRepository
}

// @inject
func NewAuthorsListUsecase(repository repository.AuthorRepository) *AuthorsListUsecase {
	return &AuthorsListUsecase{repository: repository}
}

func (this AuthorsListUsecase) Execute(page, pageSize int, name, scienceID string, journalID *uint) (*entity2.PagingEntity[entity2.AuthorEntity], error) {
	paging, err := this.repository.GetAuthorsWithArticleCount(page, pageSize, name, scienceID, journalID)
	if err != nil {
		return nil, err
	}

	return paging, nil
}
