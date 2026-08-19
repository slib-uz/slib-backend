package service

import (
	"time"

	"github.com/labstack/gommon/log"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/security"
	"slib.uz/src/infrastructure/config"
)

type AnonymousUserTokenService struct {
	tokenService security.TokenService
	env          *config.Config
}

// @inject
func NewAnonymousUserTokenService(tokenService security.TokenService, env *config.Config) *AnonymousUserTokenService {
	return &AnonymousUserTokenService{tokenService: tokenService, env: env}
}

func (this *AnonymousUserTokenService) GenerateToken(anonymousID string) string {

	accessClaim := entity.NewTokenEntity(
		this.AccessExpiredAt(),
		anonymousID,
		map[string]interface{}{
			"user_id": anonymousID,
			"type":    enum.TokenTypeAccess,
		},
	)

	return this.tokenService.Encode(accessClaim)
}

func (this *AnonymousUserTokenService) VerifyToken(tokenString string) (string, error) {
	var token, errDecode = this.tokenService.Decode(tokenString)

	if errDecode != nil {
		log.Error("AnonymousUserTokenService.VerifyToken: ", errDecode.Error())
		return "", errDecode
	}

	if token.Payload["type"] != string(enum.TokenTypeAccess) {
		return "", response.InvalidTokenError
	}

	return token.Subject, nil
}

func (this *AnonymousUserTokenService) AccessExpiredAt() time.Time {
	return time.Now().Add(time.Minute * time.Duration(this.env.ViewsCountLifetimeMinute))
}

func (this *AnonymousUserTokenService) RefreshExpiredAt() time.Time {
	return time.Now().Add(time.Hour * 24 * 30)
}
