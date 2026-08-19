package middlewares

import (
	"errors"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/logusecases"
	//"slib.uz/src/core/utils"
)

type ResponseMiddleware struct {
	adminAlert *logusecases.AdminAlertUseCase
}

// @inject
func NewResponseMiddleware(adminAlert *logusecases.AdminAlertUseCase) *ResponseMiddleware {
	return &ResponseMiddleware{adminAlert: adminAlert}
}

func (this *ResponseMiddleware) Call(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		err := next(c)

		if err != nil {

			var resp *response.Response

			if errors.As(err, &resp) {
				return c.JSON(resp.Status, resp)

			}
		}

		return err
	}
}
