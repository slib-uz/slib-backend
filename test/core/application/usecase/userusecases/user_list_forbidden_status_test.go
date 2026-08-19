package userusecases_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/userusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/entrypoint/presentation/interceptor/middlewares"
)

// fakeListRepo — UserRepository embedded, faqat GetAll kerak.
type fakeListRepo struct {
	repository.UserRepository
	calls int
}

func (f *fakeListRepo) GetAll(page, pageSize int, search string) (*entity.PagingEntity[entity.UserEntity], error) {
	f.calls++
	return &entity.PagingEntity[entity.UserEntity]{}, nil
}

// TestPlainUserGetsForbiddenNot500 — oddiy user /user/list so'raganda mijoz 403
// olishi kerak, 500 emas. Bu test butun zanjirni tekshiradi:
// usecase xatosi -> ResponseMiddleware -> Echo -> HTTP status.
func TestPlainUserGetsForbiddenNot500(t *testing.T) {
	repo := &fakeListRepo{}
	uc := userusecases.NewUserListUseCase(repo)

	plainUser := &entity.UserBasicEntity{ID: 42}

	_, err := uc.Execute(1, 10, "", plainUser)
	if err == nil {
		t.Fatal("oddiy user ro'yxatni oldi — huquq tekshiruvi ishlamadi")
	}

	// ResponseMiddleware'dagi aynan shu shart (response_middleware.go:30)
	// xatoni HTTP javobga aylantiradi.
	var resp *response.Response
	if !errors.As(err, &resp) {
		t.Fatalf("xato *response.Response emas (%T) — ResponseMiddleware uni "+
			"tanimaydi va Echo default handler 500 qaytaradi", err)
	}

	if resp.Status != http.StatusForbidden {
		t.Fatalf("status %d, kutilgani 403", resp.Status)
	}

	if repo.calls != 0 {
		t.Fatalf("rad etilgan so'rov uchun DB'ga borildi: %d", repo.calls)
	}
}

// TestForbiddenErrorProducesHTTP403 — xato haqiqiy Echo zanjiridan o'tganda
// mijozga qaysi status yetib borishini o'lchaydi.
func TestForbiddenErrorProducesHTTP403(t *testing.T) {
	repo := &fakeListRepo{}
	uc := userusecases.NewUserListUseCase(repo)

	e := echo.New()

	handler := func(c echo.Context) error {
		_, err := uc.Execute(1, 10, "", &entity.UserBasicEntity{ID: 42})
		return err
	}

	// Haqiqiy ResponseMiddleware (app.go:274 da xuddi shu tarzda ro'yxatdan
	// o'tadi). adminAlert nil beriladi: Call ichida u hech qachon
	// ishlatilmaydi (faqat konstruktorda saqlanadi), shuning uchun bu yerda
	// xavfsiz.
	responseMW := middlewares.NewResponseMiddleware(nil).Call

	e.GET("/user/list", handler, responseMW)

	req := httptest.NewRequest(http.MethodGet, "/user/list", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("mijoz %d oldi, kutilgani 403 — javob: %s", rec.Code, rec.Body.String())
	}
}
