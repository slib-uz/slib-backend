package gateway

import (
	"slib.uz/src/core/domain/entity"
)

type DegreeGateway interface {
	GetDegree(pin string) (*entity.DegreeEntity, error)
}
