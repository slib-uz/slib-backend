package gateway

import (
	"slib.uz/src/core/domain/entity"
)

type ScienceIdOAuthGateway interface {
	GetUserByCode(code string) (*entity.UserEntity, error)
	AuthorizeURL(redirectURL string) string
}
