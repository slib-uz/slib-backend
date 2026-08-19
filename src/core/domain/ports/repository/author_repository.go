package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
)

type AuthorRepository interface {
	GetByScienceID(scienceID string) (*entity2.AuthorEntity, error)
	SaveAuthor(author *entity2.AuthorEntity) (*entity2.AuthorEntity, error)

	ExistingIds(ids []uint) ([]uint, error)

	GetAuthorsWithArticleCount(page, pageSize int, name, scienceID string, journalID *uint) (*entity2.PagingEntity[entity2.AuthorEntity], error)
	GetTopAuthorsWithArticleCount(top int, journalID *uint) ([]*entity2.AuthorEntity, error)
	GetOwnCoAuthorId(scienceID string) (uint, error)

	GetJobs(ids []uint) ([]*entity2.JobWithAuthorIDEntity, error)
}
