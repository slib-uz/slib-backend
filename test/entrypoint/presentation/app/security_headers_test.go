package app_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	apppkg "slib.uz/src/entrypoint/presentation/app"
)

func TestSecurityHeadersArePresent(t *testing.T) {
	e := echo.New()
	e.Use(apppkg.SecurityHeadersMiddleware())
	e.GET("/probe", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
	})

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/probe", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("200 kutilgandi, %d keldi", rec.Code)
	}

	want := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
	}
	for header, expected := range want {
		if got := rec.Header().Get(header); got != expected {
			t.Errorf("%s: %q kutilgandi, %q keldi", header, expected, got)
		}
	}

	if csp := rec.Header().Get("Content-Security-Policy"); !strings.Contains(csp, "default-src 'none'") {
		t.Errorf("API CSP default-src 'none' bo'lishi kerak edi, %q keldi", csp)
	}
}

// Swagger UI HTML + JS SPA: default-src 'none' brauzerda oq oynaga olib keladi.
func TestSwaggerDocsCspAllowsUiAssets(t *testing.T) {
	e := echo.New()
	e.Use(apppkg.SecurityHeadersMiddleware())
	e.GET("/docs/index.html", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<html></html>")
	})

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/docs/index.html", nil))

	csp := rec.Header().Get("Content-Security-Policy")
	if strings.Contains(csp, "default-src 'none'") {
		t.Fatalf("Swagger UI default-src 'none' bilan ochilmaydi, CSP: %q", csp)
	}
	for _, token := range []string{"script-src", "style-src", "connect-src"} {
		if !strings.Contains(csp, token) {
			t.Errorf("CSP da %s bo'lishi kerak, keldi: %q", token, csp)
		}
	}
	if !strings.Contains(csp, "frame-ancestors 'none'") {
		t.Errorf("frame-ancestors saqlanishi kerak, CSP: %q", csp)
	}
}

// Eski X-XSS-Protection zamonaviy brauzerlarda o'chirilgan va yoqilganda
// o'zi zaiflik manbai bo'lishi mumkin.
func TestLegacyXssProtectionIsDisabled(t *testing.T) {
	e := echo.New()
	e.Use(apppkg.SecurityHeadersMiddleware())
	e.GET("/probe", func(c echo.Context) error { return c.NoContent(http.StatusOK) })

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/probe", nil))

	if got := rec.Header().Get("X-XSS-Protection"); got != "" && got != "0" {
		t.Errorf("X-XSS-Protection o'chirilgan bo'lishi kerak edi, %q keldi", got)
	}
}
