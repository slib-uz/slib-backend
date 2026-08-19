package authusecases_test

import (
	"context"
	"testing"
	"time"

	"slib.uz/src/core/application/service"
	"slib.uz/src/core/application/usecase/authusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/conf"
	"slib.uz/src/core/domain/ports/repository"
)

// fakeUserRepo — UserRepository interfeysi struct ichiga joylangan (embedded),
// shuning uchun faqat GetById qayta yoziladi.
type fakeUserRepo struct {
	repository.UserRepository
	user *entity.UserEntity
}

func (f *fakeUserRepo) GetById(id uint) (*entity.UserEntity, error) {
	return f.user, nil
}

// fakeConfig — ConfigAdapter interfeysi embedded; faqat kerakli getterlar yoziladi.
type fakeConfig struct {
	conf.ConfigAdapter
	strict bool
	grace  int
}

func (f *fakeConfig) IsRefreshRotationStrict() bool        { return f.strict }
func (f *fakeConfig) GetRefreshRotationGraceSeconds() int  { return f.grace }
func (f *fakeConfig) GetJwtAccessTokenExpireMinutes() int  { return 15 }
func (f *fakeConfig) GetJwtRefreshTokenExpireMinutes() int { return 10080 }

func newRefreshUseCase(c *fakeCache, strict bool, token *entity.TokenEntity) *authusecases.RefreshTokenUseCase {
	revocation := service.NewTokenRevocationService(c)
	cfg := &fakeConfig{strict: strict, grace: 60}
	auth := service.NewUserAuthTokenService(
		&fakeTokenService{token: token},
		&fakeUserRepo{user: &entity.UserEntity{ID: 42}},
		cfg,
		revocation,
	)
	return authusecases.NewRefreshTokenUseCase(auth, revocation, cfg)
}

func TestRefreshKeepsOldTokenAliveInPhaseOne(t *testing.T) {
	c := newFakeCache()
	uc := newRefreshUseCase(c, false, refreshTokenOf("42", "old-jti"))

	result, err := uc.Execute(context.Background(), "raw-refresh")
	if err != nil {
		t.Fatalf("refresh xato qaytardi: %v", err)
	}
	if result.AccessToken == "" || result.RefreshToken == "" {
		t.Fatal("yangi juftlik qaytarilmadi")
	}
	if _, revoked := c.data["revoked:jti:old-jti"]; revoked {
		t.Fatal("1-bosqichda eski refresh bekor qilindi — eski frontend buziladi")
	}
}

func TestRefreshRevokesOldTokenInPhaseTwo(t *testing.T) {
	c := newFakeCache()
	uc := newRefreshUseCase(c, true, refreshTokenOf("42", "old-jti"))

	if _, err := uc.Execute(context.Background(), "raw-refresh"); err != nil {
		t.Fatalf("refresh xato qaytardi: %v", err)
	}
	if _, revoked := c.data["revoked:jti:old-jti"]; !revoked {
		t.Fatal("2-bosqichda eski refresh bekor qilinmadi")
	}
}

// Rotatsiyaning butun maqsadi shu: 2-bosqichda aylantirib chiqarilgan refresh
// token qayta ishlatilmasligi kerak. Yuqoridagi test faqat denylist'ga yozuv
// tushganini tekshiradi — bu esa oqibatni, ya'ni ikkinchi urinish rad etilishini
// bevosita mahkamlaydi. grace=0 bo'lgani uchun bekor qilingan token so'zsiz rad etiladi.
func TestRefreshRejectsReusedTokenInPhaseTwo(t *testing.T) {
	c := newFakeCache()
	token := refreshTokenOf("42", "old-jti")

	revocation := service.NewTokenRevocationService(c)
	cfg := &fakeConfig{strict: true, grace: 0}
	auth := service.NewUserAuthTokenService(
		&fakeTokenService{token: token},
		&fakeUserRepo{user: &entity.UserEntity{ID: 42}},
		cfg,
		revocation,
	)
	uc := authusecases.NewRefreshTokenUseCase(auth, revocation, cfg)

	// Birinchi refresh muvaffaqiyatli bo'ladi va eski tokenni bekor qiladi.
	if _, err := uc.Execute(context.Background(), "raw-refresh"); err != nil {
		t.Fatalf("birinchi refresh xato qaytardi: %v", err)
	}

	// O'sha eski token bilan ikkinchi urinish — endi rad etilishi shart.
	_, err := uc.Execute(context.Background(), "raw-refresh")
	if err == nil {
		t.Fatal("aylantirib chiqarilgan refresh token qayta ishlatildi — rotatsiya himoya bermayapti")
	}
	assertStatus(t, err, 401)
}

func TestRefreshAcceptsRecentlyRevokedTokenInsideGrace(t *testing.T) {
	c := newFakeCache()
	token := refreshTokenOf("42", "old-jti")

	revocation := service.NewTokenRevocationService(c)
	// Token endigina rotatsiya tufayli bekor qilindi — grace oynasi ichida
	// hali qabul qilinishi kerak (parallel so'rovlar uzilmasligi uchun).
	if err := revocation.RevokeWithGrace(context.Background(), "old-jti", token.Exp); err != nil {
		t.Fatalf("RevokeWithGrace xato qaytardi: %v", err)
	}

	cfg := &fakeConfig{strict: true, grace: 60}
	auth := service.NewUserAuthTokenService(
		&fakeTokenService{token: token},
		&fakeUserRepo{user: &entity.UserEntity{ID: 42}},
		cfg,
		revocation,
	)
	uc := authusecases.NewRefreshTokenUseCase(auth, revocation, cfg)

	if _, err := uc.Execute(context.Background(), "raw-refresh"); err != nil {
		t.Fatalf("grace oynasi ichidagi refresh rad etildi: %v", err)
	}
}

// Grace oynasi faqat rotatsiya uchun. Logout qilingan refresh token grace ichida
// ham rad etilishi shart: aks holda o'g'irlangan refresh tokeni bo'lgan hujumchi
// jabrlanuvchi chiqqanini ko'rib, 60 soniya ichida yangi juftlik olib, 60 soniyalik
// oynani 7 kunlik yangi sessiyaga aylantirardi.
func TestRefreshRejectsLogoutRevokedTokenInsideGrace(t *testing.T) {
	c := newFakeCache()
	token := refreshTokenOf("42", "old-jti")

	revocation := service.NewTokenRevocationService(c)

	// Logout aynan shu yo'ldan yuradi: LogoutUseCase ham Revoke ni chaqiradi.
	if err := revocation.Revoke(context.Background(), "old-jti", token.Exp); err != nil {
		t.Fatalf("Revoke xato qaytardi: %v", err)
	}

	cfg := &fakeConfig{strict: true, grace: 60}
	auth := service.NewUserAuthTokenService(
		&fakeTokenService{token: token},
		&fakeUserRepo{user: &entity.UserEntity{ID: 42}},
		cfg,
		revocation,
	)
	uc := authusecases.NewRefreshTokenUseCase(auth, revocation, cfg)

	_, err := uc.Execute(context.Background(), "raw-refresh")
	if err == nil {
		t.Fatal("logout qilingan refresh token grace oynasida qabul qilindi — chiqish bekor qilinardi")
	}
	assertStatus(t, err, 401)
}

// Yuqoridagi test bilan bir xil sozlama, lekin logout haqiqiy LogoutUseCase
// orqali bajariladi: bu logout va refresh o'rtasidagi chokni mahkamlaydi.
func TestRefreshRejectsTokenAfterLogoutInsideGrace(t *testing.T) {
	c := newFakeCache()
	token := refreshTokenOf("42", "old-jti")
	tokenService := &fakeTokenService{token: token}

	revocation := service.NewTokenRevocationService(c)

	logout := authusecases.NewLogoutUseCase(revocation, tokenService)
	if err := logout.Execute(context.Background(), 42, "access-jti", time.Now().Add(15*time.Minute), "raw-refresh"); err != nil {
		t.Fatalf("logout xato qaytardi: %v", err)
	}

	cfg := &fakeConfig{strict: true, grace: 60}
	auth := service.NewUserAuthTokenService(
		tokenService,
		&fakeUserRepo{user: &entity.UserEntity{ID: 42}},
		cfg,
		revocation,
	)
	uc := authusecases.NewRefreshTokenUseCase(auth, revocation, cfg)

	_, err := uc.Execute(context.Background(), "raw-refresh")
	if err == nil {
		t.Fatal("logoutdan keyin refresh token qabul qilindi — spec 12.2 mezoni buzildi")
	}
	assertStatus(t, err, 401)
}
