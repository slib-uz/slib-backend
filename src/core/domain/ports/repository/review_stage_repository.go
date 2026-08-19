package repository

import (
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type ReviewStageRepository interface {
	GetByID(id uint) (*entity.ReviewStageEntity, error)
	GetWithApplication(id uint) (*entity.ReviewStageEntity, error)
	GetByAppID(appId uint) ([]*entity.ReviewStageEntity, error)
	Create(reviewStage *entity.ReviewStageEntity) (uint, error)
	UpdateToNext(reviewStage *entity.ReviewStageEntity, nextStage *entity.ReviewStageEntity) error
	GetBy(appId uint, stage enum.Stage, status enum.Status) (*entity.ReviewStageEntity, error)
	GetByStage(appId uint, stage enum.Stage) (*entity.ReviewStageEntity, error)
	IsExists(appId uint, stage enum.Stage, status enum.Status) (bool, error)
	GetByIDWithArticle(id uint) (*entity.ReviewStageEntity, error)
	UpdateDeadline(id uint, deadline time.Time, deadlineType enum.DeadlineType) error
	GetExpiringStages() ([]*entity.ReviewStageEntity, error)
	GetStatistics(journalID *uint) ([]*entity.ReviewStageStatisticsEntity, error)
	GetOverdueStatistics(journalID *uint) ([]*entity.ReviewStageOverdueStatisticsEntity, error)
	GetOverdueList(page, pageSize int, journalID *uint) (*entity.PagingEntity[entity.ReviewStageOverdueEntity], error)
}
