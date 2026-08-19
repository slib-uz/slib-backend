package middlewares_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/service"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	appcontext "slib.uz/src/entrypoint/presentation/app/context"
	"slib.uz/src/entrypoint/presentation/interceptor/middlewares"
)

// fakeTokenService — security.TokenService portining soxta amalga oshirilishi.
type fakeTokenService struct {
	token *entity.TokenEntity
}

func (f *fakeTokenService) Encode(claim *entity.TokenEntity) string { return "signed" }

func (f *fakeTokenService) Decode(token string) (*entity.TokenEntity, error) {
	return f.token, nil
}

// fakeUserRepo va fakeUserProfileRepo — interfeyslar embedded, faqat kerakli
// metodlar qayta yoziladi; qolganlari chaqirilsa panic beradi.
type fakeUserRepo struct {
	repository.UserRepository
}

func (f *fakeUserRepo) GetById(id uint) (*entity.UserEntity, error) {
	return &entity.UserEntity{ID: id}, nil
}

type fakeUserProfileRepo struct {
	repository.UserProfileRepository
}

func (f *fakeUserProfileRepo) UpdateLastOnlineAt(userID uint) error { return nil }

// emptyCache — cache.CacheProvider porti; denylist bo'sh, ya'ni hech bir token
// bekor qilinmagan.
type emptyCache struct{}

func newEmptyCache() *emptyCache { return &emptyCache{} }

func (f *emptyCache) GetByKey(ctx context.Context, key string) (string, error) { return "", nil }

func (f *emptyCache) Set(ctx context.Context, key string, value string, exp time.Duration) error {
	return nil
}

// Logout c.TokenID va c.TokenExp ga to'liq bog'liq: ular bo'sh kelsa, logout
// hech narsani bekor qilmay 200 qaytaradi. Ularni to'ldiradigan yagona joy shu
// middleware, va bu chok boshqa hech bir test bilan qoplanmagan.
func TestJwtAuthMiddlewarePopulatesTokenIDAndExp(t *testing.T) {
	exp := time.Now().Add(15 * time.Minute).Truncate(time.Second)

	token := entity.NewTokenEntity(exp, "42", map[string]any{
		"type": string(enum.TokenTypeAccess),
	})
	token.ID = "jti-1"

	authService := service.NewUserAuthTokenService(
		&fakeTokenService{token: token},
		&fakeUserRepo{},
		nil,
		service.NewTokenRevocationService(newEmptyCache()),
	)

	middleware := middlewares.NewJwAuthMiddleware(authService, &fakeUserProfileRepo{})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer raw-access-token")
	c := appcontext.NewContext(e.NewContext(req, httptest.NewRecorder()))

	called := false
	handler := middleware.Call(func(echo.Context) error {
		called = true
		return nil
	})

	if err := handler(c); err != nil {
		t.Fatalf("middleware xato qaytardi: %v", err)
	}
	if !called {
		t.Fatal("keyingi handler chaqirilmadi")
	}

	if c.TokenID != "jti-1" {
		t.Fatalf("c.TokenID to'ldirilmadi: %q — logout tokenni bekor qila olmasdi", c.TokenID)
	}
	if !c.TokenExp.Equal(exp) {
		t.Fatalf("c.TokenExp to'ldirilmadi: %v, kutilgan: %v", c.TokenExp, exp)
	}
	if c.User == nil || c.User.ID != 42 {
		t.Fatal("c.User to'ldirilmadi")
	}
}
