package repository

import (
	"time"

	entity2 "slib.uz/src/core/domain/entity"
)

type ApplicationRepository interface {
	Create(article *entity2.ArticleCreateEntity, userId uint, technicalReviewDeadline time.Time) (*entity2.ApplicationEntity, error)
	GetByJournalID(journalId uint, page, pageSize int, startDate, endDate time.Time, search string, status int) (*entity2.PagingEntity[entity2.ApplicationEntity], error)
	GetByIDWithRelations(id uint) (*entity2.ApplicationEntity, error)
	FindByID(id uint) (*entity2.ApplicationEntity, error)
	GetWithArticleAndJournal(id uint) (*entity2.ApplicationEntity, error)
	GetByUserID(userId uint, page, pageSize int, journalID uint) (*entity2.PagingEntity[entity2.ApplicationEntity], error)
	Update(applicationID uint, article *entity2.ArticleUpdateEntity) error
	Publish(articleID, applicationID uint, finalFilePath string) error
	GetUserAppByID(id, userID uint) (*entity2.ApplicationEntity, error)
	UpdateFile(applicationID uint, filePath string) error
	CheckUniqueConstraints(article *entity2.ArticleCreateEntity) error
	FindByIDWithRelations(id uint) (*entity2.ApplicationEntity, error)
	GetApplicationCountByUserID(userID uint, isPublished bool) (int64, error)
	GetCheckForPaymentByApplicationID(applicationID uint) (*entity2.ReviewStageEntity, error)
}
