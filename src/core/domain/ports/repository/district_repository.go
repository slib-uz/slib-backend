package repository

import (
	"slib.uz/src/core/domain/entity"
)

type DistrictRepository interface {
	ListByRegionID(regionID uint) ([]*entity.DistrictEntity, error)
}
