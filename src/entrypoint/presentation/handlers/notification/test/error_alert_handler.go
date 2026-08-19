package test

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
)

type ErrorAlertHandler struct {
}

// @inject
func NewErrorAlertHandler() *ErrorAlertHandler {
	return &ErrorAlertHandler{}
}

// Handle godoc
// @Tags notification-test
// @Accept json
// @Produce json
// @Security BasicAuth
// @Param mode query string false "Mode: 0 for panic, 1 for error"
// @Router /notification-test/error-alert [get]
func (this *ErrorAlertHandler) Handle(ctx echo.Context) error {

	mode := ctx.QueryParam("mode")
	if mode == "0" {
		panic(fmt.Errorf("ErrorAlertHandler: panic error"))
	}
	if mode == "1" {
		return response.NewFailResponse(500, "ErrorAlertHandler: 1 mode")
	}

	return response.NewFailResponse(400, "ErrorAlertHandler: invalid mode")
}
