package security

import (
	"slib.uz/src/core/domain/entity"
)

type TokenService interface {
	Encode(claim *entity.TokenEntity) string
	Decode(token string) (*entity.TokenEntity, error)
}
