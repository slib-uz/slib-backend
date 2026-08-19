package service_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"slib.uz/src/core/application/service"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

// fakeTokenService — TokenService portining soxta amalga oshirilishi.
type fakeTokenService struct {
	token *entity.TokenEntity
	err   error
}

func (f *fakeTokenService) Encode(claim *entity.TokenEntity) string { return "signed" }

func (f *fakeTokenService) Decode(token string) (*entity.TokenEntity, error) {
	return f.token, f.err
}

// fakeUserRepo — UserRepository interfeysi struct ichiga joylangan (embedded),
// shuning uchun faqat GetById qayta yoziladi; qolgan 15 metod chaqirilsa panic beradi.
type fakeUserRepo struct {
	repository.UserRepository
	user *entity.UserEntity
}

func (f *fakeUserRepo) GetById(id uint) (*entity.UserEntity, error) {
	return f.user, nil
}

func newAuthService(revocation *service.TokenRevocationService, token *entity.TokenEntity) *service.UserAuthTokenService {
	return service.NewUserAuthTokenService(
		&fakeTokenService{token: token},
		&fakeUserRepo{user: &entity.UserEntity{ID: 42}},
		nil,
		revocation,
	)
}

func accessToken(jti string, exp time.Time) *entity.TokenEntity {
	t := entity.NewTokenEntity(exp, "42", map[string]any{"type": string(enum.TokenTypeAccess)})
	t.ID = jti
	return t
}

func TestVerifyTokenAcceptsLiveToken(t *testing.T) {
	revocation := service.NewTokenRevocationService(newFakeCache())
	exp := time.Now().Add(15 * time.Minute)
	svc := newAuthService(revocation, accessToken("jti-1", exp))

	user, token, err := svc.VerifyToken(context.Background(), "raw", enum.TokenTypeAccess, 0)
	if err != nil {
		t.Fatalf("yaroqli token rad etildi: %v", err)
	}
	if user == nil || user.ID != 42 {
		t.Fatal("foydalanuvchi qaytarilmadi")
	}
	if token == nil || token.ID != "jti-1" {
		t.Fatal("dekodlangan token qaytarilmadi")
	}
}

func TestVerifyTokenRejectsRevokedToken(t *testing.T) {
	revocation := service.NewTokenRevocationService(newFakeCache())
	exp := time.Now().Add(15 * time.Minute)

	if err := revocation.Revoke(context.Background(), "jti-1", exp); err != nil {
		t.Fatalf("Revoke xato qaytardi: %v", err)
	}

	svc := newAuthService(revocation, accessToken("jti-1", exp))

	_, _, err := svc.VerifyToken(context.Background(), "raw", enum.TokenTypeAccess, 0)
	if err == nil {
		t.Fatal("bekor qilingan token qabul qilindi")
	}
}

func TestVerifyTokenAcceptsRotationRevokedTokenInsideGrace(t *testing.T) {
	revocation := service.NewTokenRevocationService(newFakeCache())
	exp := time.Now().Add(15 * time.Minute)

	if err := revocation.RevokeWithGrace(context.Background(), "jti-1", exp); err != nil {
		t.Fatalf("RevokeWithGrace xato qaytardi: %v", err)
	}

	svc := newAuthService(revocation, accessToken("jti-1", exp))

	// Endigina rotatsiya tufayli bekor qilindi, grace 60 soniya — hali qabul qilinishi kerak.
	_, _, err := svc.VerifyToken(context.Background(), "raw", enum.TokenTypeAccess, time.Minute)
	if err != nil {
		t.Fatalf("grace oynasi ichidagi token rad etildi: %v", err)
	}
}

// Logout bilan bekor qilingan token grace oynasiga bo'ysunmasligi shart.
// Aks holda o'g'irlangan refresh token bilan chiqishdan keyingi 60 soniya ichida
// yangi juftlik olib, logoutni bekor qilish mumkin bo'lardi.
func TestVerifyTokenRejectsLogoutRevokedTokenInsideGrace(t *testing.T) {
	revocation := service.NewTokenRevocationService(newFakeCache())
	exp := time.Now().Add(15 * time.Minute)

	if err := revocation.Revoke(context.Background(), "jti-1", exp); err != nil {
		t.Fatalf("Revoke xato qaytardi: %v", err)
	}

	svc := newAuthService(revocation, accessToken("jti-1", exp))

	_, _, err := svc.VerifyToken(context.Background(), "raw", enum.TokenTypeAccess, time.Minute)
	if err == nil {
		t.Fatal("logout bilan bekor qilingan token grace oynasida qabul qilindi")
	}
}

func TestVerifyTokenRejectsWrongTokenType(t *testing.T) {
	revocation := service.NewTokenRevocationService(newFakeCache())
	svc := newAuthService(revocation, accessToken("jti-1", time.Now().Add(time.Minute)))

	_, _, err := svc.VerifyToken(context.Background(), "raw", enum.TokenTypeRefresh, 0)
	if err == nil {
		t.Fatal("access token refresh sifatida qabul qilindi")
	}
}

// countingUserRepo GetById necha marta chaqirilganini sanaydi.
type countingUserRepo struct {
	repository.UserRepository
	user  *entity.UserEntity
	calls int
}

func (f *countingUserRepo) GetById(id uint) (*entity.UserEntity, error) {
	f.calls++
	return f.user, nil
}

// Denylist tekshiruvi DB'dan oldin turishi kerak: bekor qilingan token uchun
// ortiqcha DB so'rovi ketmasligi lozim.
func TestVerifyTokenSkipsDatabaseWhenTokenRevoked(t *testing.T) {
	revocation := service.NewTokenRevocationService(newFakeCache())
	exp := time.Now().Add(15 * time.Minute)

	if err := revocation.Revoke(context.Background(), "jti-1", exp); err != nil {
		t.Fatalf("Revoke xato qaytardi: %v", err)
	}

	repo := &countingUserRepo{user: &entity.UserEntity{ID: 42}}
	svc := service.NewUserAuthTokenService(
		&fakeTokenService{token: accessToken("jti-1", exp)},
		repo,
		nil,
		revocation,
	)

	if _, _, err := svc.VerifyToken(context.Background(), "raw", enum.TokenTypeAccess, 0); err == nil {
		t.Fatal("bekor qilingan token qabul qilindi")
	}
	if repo.calls != 0 {
		t.Fatalf("bekor qilingan token uchun DB'ga %d marta borildi, 0 kutilgandi", repo.calls)
	}
}

// grace == 0 da rad etish so'zsiz bo'lishi kerak — soatga bog'liq emas.
// Bu yerda bekor qilish vaqti tekshiruvchi soatdan oldinda: instancelar orasidagi
// manfiy clock skew ni taqlid qiladi. time.Since(*revokedAt) manfiy chiqadi,
// shuning uchun ">" solishtiruvi yolg'iz o'zi tokenni o'tkazib yuborardi.
func TestVerifyTokenRejectsRevokedTokenWithZeroGraceDespiteClockSkew(t *testing.T) {
	cache := newFakeCache()
	exp := time.Now().Add(15 * time.Minute)

	// Bekor qiluvchi instance soati 2 soniya oldinda.
	future := time.Now().Add(2 * time.Second).Unix()
	cache.data[service.RevokedKeyPrefix+"jti-1"] = strconv.FormatInt(future, 10)

	svc := newAuthService(service.NewTokenRevocationService(cache), accessToken("jti-1", exp))

	_, _, err := svc.VerifyToken(context.Background(), "raw", enum.TokenTypeAccess, 0)
	if err == nil {
		t.Fatal("grace=0 bo'lsa ham bekor qilingan token qabul qilindi (clock skew)")
	}
}

// Task 1 dan oldin berilgan tokenlarda jti yo'q (ID == ""). Ular oddiy tarzda
// tekshirilishi shart: aks holda deploy paytida hamma mavjud sessiya uziladi.
func TestVerifyTokenAcceptsLegacyTokenWithoutJti(t *testing.T) {
	revocation := service.NewTokenRevocationService(newFakeCache())
	exp := time.Now().Add(15 * time.Minute)
	svc := newAuthService(revocation, accessToken("", exp))

	user, token, err := svc.VerifyToken(context.Background(), "raw", enum.TokenTypeAccess, 0)
	if err != nil {
		t.Fatalf("jti siz eski token rad etildi: %v", err)
	}
	if user == nil || user.ID != 42 {
		t.Fatal("foydalanuvchi qaytarilmadi")
	}
	if token == nil || token.ID != "" {
		t.Fatal("dekodlangan token qaytarilmadi")
	}
}

// Redis yiqilganda auth yo'li fail-OPEN bo'lishi kerak: denylist o'qib bo'lmasa,
// foydalanuvchi ichkariga kiritiladi. Bu ongli qaror — xavf oynasi access token
// TTL si bilan chegaralangan, ammo Redis nosozligi butun saytni to'xtatmaydi.
func TestVerifyTokenFailsOpenWhenCacheUnavailable(t *testing.T) {
	cache := newFakeCache()
	cache.getErr = errors.New("redis down")
	exp := time.Now().Add(15 * time.Minute)

	svc := newAuthService(service.NewTokenRevocationService(cache), accessToken("jti-1", exp))

	user, _, err := svc.VerifyToken(context.Background(), "raw", enum.TokenTypeAccess, 0)
	if err != nil {
		t.Fatalf("cache yiqilganda token rad etildi, fail-open kutilgandi: %v", err)
	}
	if user == nil || user.ID != 42 {
		t.Fatal("foydalanuvchi qaytarilmadi")
	}
}
