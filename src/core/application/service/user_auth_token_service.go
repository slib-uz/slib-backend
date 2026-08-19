package service

import (
	"context"
	"strconv"
	"time"

	"github.com/labstack/gommon/log"
	"slib.uz/src/core/application/mapper"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/conf"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/security"
)

type UserAuthTokenService struct {
	tokenService   security.TokenService
	userRepository repository.UserRepository
	cfg            conf.ConfigAdapter
	revocation     *TokenRevocationService
}

// @inject
func NewUserAuthTokenService(
	tokenService security.TokenService,
	userRepository repository.UserRepository,
	cfg conf.ConfigAdapter,
	revocation *TokenRevocationService,
) *UserAuthTokenService {
	return &UserAuthTokenService{
		tokenService:   tokenService,
		userRepository: userRepository,
		cfg:            cfg,
		revocation:     revocation,
	}
}

func (this *UserAuthTokenService) GenerateToken(userID uint) *entity.AuthTokenEntity {
	return entity.NewAuthTokenEntity(this.GenerateAccessToken(userID), this.GenerateRefreshToken(userID))
}

func (this *UserAuthTokenService) GenerateAccessToken(userID uint) string {
	user, err := this.userRepository.GetById(userID)
	if err != nil {
		return ""
	}

	type role struct {
		Role          enum.UserRole `json:"role"`
		JournalID     *uint         `json:"journal_id,omitempty"`
		PublisherID   *uint         `json:"publisher_id,omitempty"`
		InstitutionID *uint         `json:"institution_id,omitempty"`
	}

	roles := make([]*role, len(user.Roles))

	for i, it := range user.Roles {
		roles[i] = &role{
			Role:          it.Role,
			PublisherID:   it.PublisherID,
			JournalID:     it.JournalID,
			InstitutionID: it.InstitutionID,
		}
	}

	accessClaim := entity.NewTokenEntity(
		this.GetAccessExpiredAt(),
		strconv.Itoa(int(userID)),
		map[string]any{
			"type":         enum.TokenTypeAccess,
			"user_id":      userID,
			"science_id":   user.ScienceID,
			"full_name":    user.FullName,
			"phone_number": user.PhoneNumber,
			"roles":        roles,
		},
	)

	return this.tokenService.Encode(accessClaim)
}

func (this *UserAuthTokenService) GenerateRefreshToken(userID uint) string {
	refreshClaim := entity.NewTokenEntity(
		this.GetRefreshExpiredAt(),
		strconv.Itoa(int(userID)),
		map[string]any{
			"user_id": userID,
			"type":    enum.TokenTypeRefresh,
		},
	)
	return this.tokenService.Encode(refreshClaim)
}

// VerifyToken tokenni tekshiradi va foydalanuvchi bilan birga dekodlangan tokenni qaytaradi.
//
// grace > 0 bo'lsa, shu muddat ichida rotatsiya tufayli bekor qilingan token hali
// ham qabul qilinadi. Bu refresh rotatsiyasida parallel so'rovlar uzilmasligi uchun
// kerak. Logout bilan bekor qilingan token (Revocation.Final) grace ga bo'ysunmaydi.
// Access yo'lida grace har doim 0 bo'ladi va bekor qilingan token so'zsiz rad etiladi:
// instancelar orasidagi soat farqi (clock skew) bu qaroga ta'sir qilmasligi shart.
//
// Denylist o'qib bo'lmasa RevokedAt fail-open qiladi ((nil, nil) qaytaradi) —
// Redis yiqilishi butun saytni to'xtatmasligi kerak. Quyidagi err shoxi shu sababli
// amalda ishlamaydi, lekin RevokedAt kelajakda xato qaytara boshlasa,
// auth yo'li jimgina fail-closed'ga o'tib ketmasligi uchun saqlab qolingan.
func (this *UserAuthTokenService) VerifyToken(
	ctx context.Context,
	tokenString string,
	tokenType enum.TokenType,
	grace time.Duration,
) (*entity.UserBasicEntity, *entity.TokenEntity, error) {

	token, errDecode := this.tokenService.Decode(tokenString)
	if errDecode != nil {
		return nil, nil, errDecode
	}

	if token.Payload["type"] != string(tokenType) {
		return nil, nil, response.InvalidTokenError
	}

	revoked, err := this.revocation.RevokedAt(ctx, token.ID)
	if err != nil {
		return nil, nil, err
	}
	// grace <= 0 — so'zsiz rad etish. Buni time.Since dan kelib chiqarib bo'lmaydi:
	// tekshiruvchi instance soati bekor qiluvchisidan orqada bo'lsa, farq davomida
	// bekor qilingan token yaroqli ko'rinib qolardi.
	//
	// revoked.Final — logout bilan bekor qilingan token: grace unga umuman
	// tegishli emas. Aks holda logoutdan keyingi grace oynasi ichida o'g'irlangan
	// refresh token yangi juftlik olib, chiqishni bekor qila olardi.
	if revoked != nil && (revoked.Final || grace <= 0 || time.Since(revoked.At) > grace) {
		log.Warn("UserAuthTokenService.VerifyToken: bekor qilingan token ishlatildi. jti=", token.ID, " subject=", token.Subject)
		return nil, nil, response.InvalidTokenError
	}

	userID, err := strconv.Atoi(token.Subject)
	if err != nil {
		return nil, nil, response.InvalidTokenError
	}

	user, err := this.userRepository.GetById(uint(userID))
	if err != nil {
		return nil, nil, err
	}

	return mapper.UserEntityToBasic(user), token, nil
}

func (this *UserAuthTokenService) GetAccessExpiredAt() time.Time {
	// return time.Now().Add(time.Hour * 24)
	return time.Now().Add(time.Minute * time.Duration(this.cfg.GetJwtAccessTokenExpireMinutes()))
}

func (this *UserAuthTokenService) GetRefreshExpiredAt() time.Time {
	// return time.Now().Add(time.Hour * 24 * 30)
	return time.Now().Add(time.Minute * time.Duration(this.cfg.GetJwtRefreshTokenExpireMinutes()))
}
