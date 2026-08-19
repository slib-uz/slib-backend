package authusecases

import (
	"context"
	"strconv"
	"time"

	"github.com/labstack/gommon/log"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/security"
)

type LogoutUseCase struct {
	revocation   *service.TokenRevocationService
	tokenService security.TokenService
}

// @inject
func NewLogoutUseCase(revocation *service.TokenRevocationService, tokenService security.TokenService) *LogoutUseCase {
	return &LogoutUseCase{revocation: revocation, tokenService: tokenService}
}

// Execute joriy access tokenni va berilgan bo'lsa refresh tokenni bekor qiladi.
//
// Faqat access tokenni bekor qilish yetarli emas: refresh qo'lda qolsa,
// mijoz darhol yangi access token oladi va logout ma'nosiz bo'ladi.
func (this *LogoutUseCase) Execute(
	ctx context.Context,
	userID uint,
	accessJTI string,
	accessExp time.Time,
	refreshToken string,
) error {

	// jti bo'sh bo'lsa Revoke ataylab hech narsa qilmaydi (aks holda barcha
	// jti siz tokenlar bitta yalang'och kalitni bo'lishardi). Lekin bu holat
	// sodir bo'lmasligi kerak: JwtTokenService.Encode har doim uuid jti qo'yadi,
	// middleware esa uni c.TokenID ga yozadi. Jimgina 200 qaytarsak, logout
	// hech narsani bekor qilmagani bilinmay qolardi — shuning uchun ERROR.
	if accessJTI == "" {
		log.Error("LogoutUseCase: access token jti si bo'sh, access bekor qilinmadi. user_id=", userID)
	}

	if err := this.revocation.Revoke(ctx, accessJTI, accessExp); err != nil {
		log.Error("LogoutUseCase: access token bekor qilinmadi. user_id=", userID, " err=", err.Error())
		return response.LogoutFailedError
	}

	if refreshToken == "" {
		log.Warn("LogoutUseCase: refresh token yuborilmadi, faqat access bekor qilindi. user_id=", userID)
		return nil
	}

	token, err := this.tokenService.Decode(refreshToken)
	if err != nil {
		// Yaroqsiz yoki muddati o'tgan refresh — chiqish baribir muvaffaqiyatli hisoblanadi.
		return nil
	}

	if token.Payload["type"] != string(enum.TokenTypeRefresh) {
		return nil
	}

	// Egalik tekshiruvi: usiz istalgan foydalanuvchi boshqasining refresh
	// tokenini bekor qilib, uni tizimdan uzib qo'ya olardi.
	if token.Subject != strconv.Itoa(int(userID)) {
		log.Warn("LogoutUseCase: begona refresh tokenni bekor qilishga urinish. user_id=", userID, " subject=", token.Subject)
		return response.PermissionDeniedError
	}

	if err := this.revocation.Revoke(ctx, token.ID, token.Exp); err != nil {
		log.Error("LogoutUseCase: refresh token bekor qilinmadi. user_id=", userID, " err=", err.Error())
		return response.LogoutFailedError
	}

	return nil
}
