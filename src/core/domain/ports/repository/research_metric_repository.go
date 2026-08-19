package repository

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type ResearchMetricRepository interface {
	UpdateOrCreate(userID uint, source enum.ResearchMetricEnum, entity *entity.ResearchMetricEntity) error
	GetByUserID(userID uint) ([]*entity.ResearchMetricEntity, error)
	GetByUserIDs(userIDs []uint) (map[uint][]*entity.ResearchMetricEntity, error)
	DeleteByIDAndUserID(id, userID uint) error
}
