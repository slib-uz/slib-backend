package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"slib.uz/src/core/application/service"
)

type DeveloperUserBasicAuthMiddleware struct {
	service *service.DeveloperUserAuthService
}

// @inject
func NewDeveloperUserBasicAuthMiddleware(service *service.DeveloperUserAuthService) *DeveloperUserBasicAuthMiddleware {
	return &DeveloperUserBasicAuthMiddleware{service: service}
}

func (this *DeveloperUserBasicAuthMiddleware) Wrap(next echo.HandlerFunc) echo.HandlerFunc {
	return middleware.BasicAuth(func(username string, password string, c echo.Context) (bool, error) {
		_, err := this.service.Authenticate(username, password)
		if err != nil {
			return false, nil
		}
		return true, nil
	})(next)
}
