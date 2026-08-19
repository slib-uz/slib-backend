package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
)

type ArticleRepository interface {
	GetByIDWithJournal(id uint) (*entity2.ArticleBasicEntity, error)
	FindByID(id uint) (*entity2.ArticleEntity, error)
	FindByIDWithAuthors(id uint) (*entity2.ArticleEntity, error)
	FindByIDsWithAuthors(ids []uint) ([]*entity2.ArticleEntity, error)
	UpdateViewsCount(map[uint]int64) error
	GetArticleCount() (int64, error)
	GetAuthorsCount() (int64, error)
	GetArticleCountByUserID(userID uint, isPublished bool) (int64, error)
	GetArticleCountByEditionID(editionID uint) (int64, error)
	UpdateROI(id, externalID uint, roi *string) error
	PublishedIDsWithoutROI() ([]uint, error)
	UpdateDOI(id uint, doi string) error
	UpdateStatus(id uint, IsSent bool) error
	ManageUpdate(id uint, e *entity2.ArticleManageUpdateEntity) error
	AddCoAuthor(articleID uint, authorID uint) error
	OnlyArticleNameByLegacyAuthorIDs(ids []uint, page, pageSize int) (*entity2.PagingEntity[entity2.ArticleOnlyNameEntity], error)
	CreatePublishArticle(article *entity2.ArticlePublishEntity, journalID uint) (uint, error)
	CreatePublishArticleWithAffiliations(article *entity2.ArticlePublishEntity, journalID uint) (uint, error)
	CheckDOIExists(doi string) (bool, error)
	CheckROIExists(roi string) (bool, error)
	CheckExternalIDExists(externalID uint) (bool, error)
}
