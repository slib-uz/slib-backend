package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	myvalidator "slib.uz/src/entrypoint/presentation/interceptor/validator"
)

// @inject
func NewEcho() *echo.Echo {
	e := echo.New()
	e.Validator = myvalidator.NewRequestValidator(validator.New())
	return e
}
