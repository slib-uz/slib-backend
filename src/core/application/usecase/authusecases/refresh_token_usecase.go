package authusecases

import (
	"context"
	"time"

	"github.com/labstack/gommon/log"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/conf"
)

// RefreshTokenResult — rotatsiyadan keyingi yangi juftlik.
type RefreshTokenResult struct {
	AccessToken  string
	RefreshToken string
}

type RefreshTokenUseCase struct {
	service    *service.UserAuthTokenService
	revocation *service.TokenRevocationService
	cfg        conf.ConfigAdapter
}

// @inject
func NewRefreshTokenUseCase(
	svc *service.UserAuthTokenService,
	revocation *service.TokenRevocationService,
	cfg conf.ConfigAdapter,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{service: svc, revocation: revocation, cfg: cfg}
}

// Execute refresh token asosida yangi access va refresh juftligini qaytaradi.
//
// 1-bosqich (REFRESH_ROTATION_STRICT=false): eski refresh bekor qilinmaydi,
// shuning uchun yangi juftlikni saqlamaydigan eski frontend ishlashda davom etadi.
// Logout baribir ishlaydi, chunki u refresh jti sini bevosita denylist'ga yozadi.
//
// 2-bosqich (true): eski refresh bekor qilinadi, grace oynasi parallel
// so'rovlarni uzilishdan saqlaydi.
func (this *RefreshTokenUseCase) Execute(ctx context.Context, refreshToken string) (*RefreshTokenResult, error) {
	strict := this.cfg.IsRefreshRotationStrict()

	grace := time.Duration(0)
	if strict {
		grace = time.Duration(this.cfg.GetRefreshRotationGraceSeconds()) * time.Second
	}

	user, token, err := this.service.VerifyToken(ctx, refreshToken, enum.TokenTypeRefresh, grace)
	if err != nil {
		return nil, err
	}

	result := &RefreshTokenResult{
		AccessToken:  this.service.GenerateAccessToken(user.ID),
		RefreshToken: this.service.GenerateRefreshToken(user.ID),
	}

	if strict {
		// RevokeWithGrace, Revoke emas: rotatsiya tufayli bekor qilingan token
		// grace oynasi ichida qabul qilinaverishi kerak. Logout esa Revoke
		// ishlatadi va u yerda grace umuman qo'llanilmaydi.
		if err := this.revocation.RevokeWithGrace(ctx, token.ID, token.Exp); err != nil {
			// Yangi juftlik allaqachon berildi — so'rovni yiqitmaymiz, lekin belgilab qo'yamiz.
			log.Error("RefreshTokenUseCase: eski refresh bekor qilinmadi. jti=", token.ID, " err=", err.Error())
		}
	}

	return result, nil
}
