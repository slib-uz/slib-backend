package authusecases_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/labstack/gommon/log"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/application/usecase/authusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

// assertStatus xatoning aniq HTTP statusini tekshiradi.
//
// "xato qaytdi" ni tekshirish yetarli emas: bu kod bazasida xato qiymati
// 200 statusini ham olib yurishi mumkin (masalan response.ConflictError),
// va middleware resp.Status nimani aytsa o'shani qaytaradi. Ya'ni statusni
// tekshirmasak, 403 ni 401 ga yoki 503 ni 200 ga aylantirib yuborish
// testlarni yashil holda qoldiradi.
func assertStatus(t *testing.T, err error, want int) {
	t.Helper()

	resp, ok := err.(*response.Response)
	if !ok || resp.Status != want {
		t.Fatalf("%d status kutilgan edi, keldi: %v", want, err)
	}
}

type fakeCache struct {
	data   map[string]string
	setErr error
	// failOnKey bo'sh bo'lsa setErr barcha yozuvlarni yiqitadi. To'ldirilsa —
	// faqat o'sha kalit yiqiladi, bu access va refresh bekor qilish
	// shoxobchalarini alohida sinash imkonini beradi.
	failOnKey string
}

func newFakeCache() *fakeCache {
	return &fakeCache{data: make(map[string]string)}
}

func (f *fakeCache) GetByKey(ctx context.Context, key string) (string, error) {
	return f.data[key], nil
}

func (f *fakeCache) Set(ctx context.Context, key string, value string, exp time.Duration) error {
	if f.setErr != nil && (f.failOnKey == "" || f.failOnKey == key) {
		return f.setErr
	}
	f.data[key] = value
	return nil
}

type fakeTokenService struct {
	token *entity.TokenEntity
	err   error
}

func (f *fakeTokenService) Encode(claim *entity.TokenEntity) string { return "signed" }

func (f *fakeTokenService) Decode(token string) (*entity.TokenEntity, error) {
	return f.token, f.err
}

func refreshTokenOf(subject string, jti string) *entity.TokenEntity {
	t := entity.NewTokenEntity(time.Now().Add(time.Hour), subject, map[string]any{
		"type": string(enum.TokenTypeRefresh),
	})
	t.ID = jti
	return t
}

func TestLogoutRevokesAccessAndRefresh(t *testing.T) {
	c := newFakeCache()
	uc := authusecases.NewLogoutUseCase(service.NewTokenRevocationService(c), &fakeTokenService{token: refreshTokenOf("42", "refresh-jti")})

	err := uc.Execute(context.Background(), 42, "access-jti", time.Now().Add(15*time.Minute), "raw-refresh")
	if err != nil {
		t.Fatalf("logout xato qaytardi: %v", err)
	}

	if _, ok := c.data["revoked:jti:access-jti"]; !ok {
		t.Error("access token bekor qilinmadi")
	}
	if _, ok := c.data["revoked:jti:refresh-jti"]; !ok {
		t.Error("refresh token bekor qilinmadi")
	}

	// Logout idempotent: allaqachon bekor qilingan tokenlar bilan qayta
	// chaqirilsa ham 200 qaytadi.
	if err := uc.Execute(context.Background(), 42, "access-jti", time.Now().Add(15*time.Minute), "raw-refresh"); err != nil {
		t.Fatalf("takroriy logout xato qaytardi: %v", err)
	}
}

func TestLogoutRejectsRefreshTokenOfAnotherUser(t *testing.T) {
	c := newFakeCache()
	// Refresh token 99-foydalanuvchiniki, lekin so'rovni 42 yubormoqda.
	uc := authusecases.NewLogoutUseCase(service.NewTokenRevocationService(c), &fakeTokenService{token: refreshTokenOf("99", "refresh-jti")})

	err := uc.Execute(context.Background(), 42, "access-jti", time.Now().Add(15*time.Minute), "raw-refresh")
	if err == nil {
		t.Fatal("begona refresh token qabul qilindi — bekor qilish orqali DoS mumkin")
	}
	assertStatus(t, err, 403)

	if _, ok := c.data["revoked:jti:refresh-jti"]; ok {
		t.Error("begona refresh token bekor qilindi")
	}
}

func TestLogoutSucceedsWithoutRefreshToken(t *testing.T) {
	c := newFakeCache()
	uc := authusecases.NewLogoutUseCase(service.NewTokenRevocationService(c), &fakeTokenService{})

	err := uc.Execute(context.Background(), 42, "access-jti", time.Now().Add(15*time.Minute), "")
	if err != nil {
		t.Fatalf("refreshsiz logout xato qaytardi: %v", err)
	}
	if _, ok := c.data["revoked:jti:access-jti"]; !ok {
		t.Error("access token bekor qilinmadi")
	}
}

func TestLogoutSucceedsWhenRefreshTokenIsInvalid(t *testing.T) {
	c := newFakeCache()
	uc := authusecases.NewLogoutUseCase(service.NewTokenRevocationService(c), &fakeTokenService{err: errors.New("invalid token")})

	err := uc.Execute(context.Background(), 42, "access-jti", time.Now().Add(15*time.Minute), "buzuq")
	if err != nil {
		t.Fatalf("yaroqsiz refresh logoutni yiqitdi: %v", err)
	}
}

// captureLogOutput global gommon logger chiqishini ushlaydi va tiklaydi.
// Core qatlami faqat gommon/log dan foydalanadi, shuning uchun boshqa yo'l yo'q.
func captureLogOutput(t *testing.T) *bytes.Buffer {
	t.Helper()

	buf := &bytes.Buffer{}
	prev := log.Output()
	log.SetOutput(buf)
	t.Cleanup(func() { log.SetOutput(prev) })

	return buf
}

// jti bo'sh bo'lsa Revoke ataylab hech narsa qilmaydi, ya'ni logout hech narsani
// bekor qilmay 200 qaytaradi. Bu "bo'lishi mumkin emas" holat — uni ushlab
// turgan yagona narsa middleware c.TokenID ni to'ldirishi. O'sha chok buzilsa,
// logout jimgina ochiq qolib ketmasligi kerak: ERROR log qoldiriladi.
func TestLogoutLogsErrorWhenAccessJtiIsEmpty(t *testing.T) {
	logs := captureLogOutput(t)

	c := newFakeCache()
	uc := authusecases.NewLogoutUseCase(service.NewTokenRevocationService(c), &fakeTokenService{})

	err := uc.Execute(context.Background(), 42, "", time.Now().Add(15*time.Minute), "")
	if err != nil {
		t.Fatalf("logout xato qaytardi: %v", err)
	}

	// Yalang'och prefiks kaliti yaratilmasligi shart — aks holda bitta logout
	// jti siz barcha tokenlarni o'chirgan bo'lardi.
	if _, ok := c.data["revoked:jti:"]; ok {
		t.Fatal("bo'sh jti uchun yalang'och prefiks kaliti yozildi")
	}

	out := logs.String()
	if !strings.Contains(out, "ERROR") || !strings.Contains(out, "jti") {
		t.Fatalf("bo'sh jti jimgina o'tib ketdi, ERROR log kutilgandi. log=%q", out)
	}
}

func TestLogoutFailsWhenCacheWriteFails(t *testing.T) {
	c := newFakeCache()
	c.setErr = errors.New("redis down")
	uc := authusecases.NewLogoutUseCase(service.NewTokenRevocationService(c), &fakeTokenService{})

	err := uc.Execute(context.Background(), 42, "access-jti", time.Now().Add(15*time.Minute), "")
	if err == nil {
		t.Fatal("Redis nosozligi yashirildi — foydalanuvchi chiqdim deb o'ylaydi, token esa ishlaydi")
	}
	assertStatus(t, err, 503)
}

// Refresh tokenni bekor qilishdagi nosozlik alohida shoxobcha: access
// muvaffaqiyatli bekor qilinib, refresh qolib ketsa ham mijozga 200 aytish
// mumkin emas — u qo'lidagi refresh bilan darhol yangi access oladi.
func TestLogoutFailsWhenRefreshRevocationFails(t *testing.T) {
	c := newFakeCache()
	c.setErr = errors.New("redis down")
	c.failOnKey = "revoked:jti:refresh-jti"
	uc := authusecases.NewLogoutUseCase(service.NewTokenRevocationService(c), &fakeTokenService{token: refreshTokenOf("42", "refresh-jti")})

	err := uc.Execute(context.Background(), 42, "access-jti", time.Now().Add(15*time.Minute), "raw-refresh")
	if err == nil {
		t.Fatal("refresh bekor qilinmagani yashirildi — mijoz yangi access token ola oladi")
	}
	assertStatus(t, err, 503)

	// Access bekor qilingani aniq — demak test aynan refresh shoxobchasini bosdi.
	if _, ok := c.data["revoked:jti:access-jti"]; !ok {
		t.Error("access token bekor qilinmadi")
	}
}
