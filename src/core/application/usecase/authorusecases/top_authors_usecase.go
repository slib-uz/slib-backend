package authorusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type TopAuthorsListUsecase struct {
	repository repository.AuthorRepository
}

// @inject
func NewTopAuthorsListUsecase(repository repository.AuthorRepository) *TopAuthorsListUsecase {
	return &TopAuthorsListUsecase{repository: repository}
}

func (this TopAuthorsListUsecase) Execute(top int, journalID *uint) ([]*entity.AuthorEntity, error) {
	authors, err := this.repository.GetTopAuthorsWithArticleCount(top, journalID)
	if err != nil {
		return nil, err
	}
	return authors, nil
}
