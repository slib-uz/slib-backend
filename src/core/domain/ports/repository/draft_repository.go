package repository

import (
	"slib.uz/src/core/domain/entity"
)

type DraftRepository interface {
	Save(draft *entity.DraftEntity) error
	GetByKey(key string) (*entity.DraftEntity, error)
}
