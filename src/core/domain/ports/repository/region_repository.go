package repository

import (
	"slib.uz/src/core/domain/entity"
)

type RegionRepository interface {
	List() ([]*entity.RegionEntity, error)
}
