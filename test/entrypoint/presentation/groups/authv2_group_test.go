package groups_test

import (
	"testing"

	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/groups"
	"slib.uz/src/entrypoint/presentation/handlers/authv2"
)

// fakeConfigAdapter is a minimal conf.ConfigAdapter stub used only to control
// IsProduction() for this test; every other method is unused by RegisterRoutes.
type fakeConfigAdapter struct {
	isProduction bool
}

func (f *fakeConfigAdapter) GetReviewDeadlineDays() int           { return 0 }
func (f *fakeConfigAdapter) GetFrontendURL() string               { return "" }
func (f *fakeConfigAdapter) GetROIFrontendURL() string            { return "" }
func (f *fakeConfigAdapter) OtpTTLMinutes() int                   { return 0 }
func (f *fakeConfigAdapter) OtpVerifyMaxAttempts() int            { return 0 }
func (f *fakeConfigAdapter) OtpSendPerMinute() int                { return 0 }
func (f *fakeConfigAdapter) OtpSendPerHour() int                  { return 0 }
func (f *fakeConfigAdapter) GetJwtAccessTokenExpireMinutes() int  { return 0 }
func (f *fakeConfigAdapter) GetJwtRefreshTokenExpireMinutes() int { return 0 }
func (f *fakeConfigAdapter) GetCrossRefSenderEmail() string       { return "" }
func (f *fakeConfigAdapter) GetClientBasicAuthCredentials() (string, string) {
	return "", ""
}
func (f *fakeConfigAdapter) IsRefreshRotationStrict() bool       { return false }
func (f *fakeConfigAdapter) GetRefreshRotationGraceSeconds() int { return 0 }
func (f *fakeConfigAdapter) IsProduction() bool                  { return f.isProduction }

func hasRoute(routes []*echo.Route, method, path string) bool {
	for _, r := range routes {
		if r.Method == method && r.Path == path {
			return true
		}
	}
	return false
}

// Prefiks prod bilan bir xil: app.go guruhni "/auth-v2" ga ulaydi
// (app/app.go:318). Testda boshqa prefiks turgani o'quvchini chalg'itardi.
func TestAuthV2Group_SandboxLoginRoute_NotRegisteredInProduction(t *testing.T) {
	e := echo.New()
	group := e.Group("/auth-v2")

	authV2Group := groups.NewAuthV2Group(
		&authv2.SendOtpHandler{},
		&authv2.VerifyAndLoginHandler{},
		&authv2.CheckPhoneNumberHandler{},
		&authv2.SandboxLoginHandler{},
		&fakeConfigAdapter{isProduction: true},
	)
	authV2Group.RegisterRoutes(group)

	if hasRoute(e.Routes(), echo.POST, "/auth-v2/sandbox/login") {
		t.Fatal("expected sandbox/login route to be absent when IsProduction() is true, but it was registered")
	}
}

func TestAuthV2Group_SandboxLoginRoute_RegisteredOutsideProduction(t *testing.T) {
	e := echo.New()
	group := e.Group("/auth-v2")

	authV2Group := groups.NewAuthV2Group(
		&authv2.SendOtpHandler{},
		&authv2.VerifyAndLoginHandler{},
		&authv2.CheckPhoneNumberHandler{},
		&authv2.SandboxLoginHandler{},
		&fakeConfigAdapter{isProduction: false},
	)
	authV2Group.RegisterRoutes(group)

	if !hasRoute(e.Routes(), echo.POST, "/auth-v2/sandbox/login") {
		t.Fatal("expected sandbox/login route to be present when IsProduction() is false, but it was not registered")
	}
}
