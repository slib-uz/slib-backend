package gateway

import (
	"slib.uz/src/core/domain/entity"
)

type DoctorateGateway interface {
	GetDoctorates(pin string) ([]*entity.DoctorateEntity, error)
}
