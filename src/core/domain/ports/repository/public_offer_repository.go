package repository

import (
	"slib.uz/src/core/domain/entity"
)

type PublicOfferRepository interface {
	Get() (*entity.PublicOfferEntity, error)
}
