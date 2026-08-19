package repository

import (
	"slib.uz/src/core/domain/entity"
)

type AiDetectRepository interface {
	Create(entity *entity.AiDetectResultEntity) error
	GetByExternalID(externalID uint) (*entity.AiDetectResultEntity, error)
	GetByID(id uint) (*entity.AiDetectResultEntity, error)
	GetByIdWithApplication(id uint) (*entity.AiDetectResultEntity, error)
	Update(entity *entity.AiDetectResultEntity) error
	FindIDByExternalID(externalID uint) (uint, error)
	GetAllByReviewStageID(reviewStageID uint) ([]*entity.AiDetectResultEntity, error)
	GetByJournalID(journalID uint, page, pageSize int) (*entity.PagingEntity[entity.AiDetectResultEntity], error)
}
