package repository

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type PublishedArticleRepository interface {
	GetByIDWithRelations(id uint) (*entity.ArticleEntity, error)
	GetByAuthor(scienceID string, page, pageSize int, search string, studyFieldIds []uint, year int, journalID *uint) (*entity.PagingEntity[entity.ArticleEntity], error)
	GetByJournal(journalID uint, page, pageSize int, search string, authorSearch, tag string, studyFieldIds []uint, year int) (*entity.PagingEntity[entity.ArticleEntity], error)
	GetAll(
		// filters
		fromYear, toYear *int, journal uint, languages []uint,
		studyfields []uint, tags []string,
		// search parameters
		articleTitle, authorName, journalName, issnPaper, issnOnline *string,
		// sort: , views_count, rating_sum, publication_date
		sort string,
		page, pageSize int,
		accessType enum.AccessType,
		hasDoi *bool,
	) (*entity.PagingEntity[entity.ArticleEntity], error)
	GetTopArticles(page, pageSize int) (*entity.PagingEntity[entity.ArticleEntity], error)
	GetByEdition(journalID, editionID uint, unassigned bool, page, pageSize int) (*entity.PagingEntity[entity.ArticleEntity], error)
}
