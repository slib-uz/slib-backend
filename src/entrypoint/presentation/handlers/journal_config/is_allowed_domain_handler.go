package journal_config

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/usecase/journalconfigusecases"
)

type IsAllowedDomainHandler struct {
	uc *journalconfigusecases.IsAllowedDomainUseCase
}

// @inject
func NewIsAllowedDomainHandler(uc *journalconfigusecases.IsAllowedDomainUseCase) *IsAllowedDomainHandler {
	return &IsAllowedDomainHandler{uc: uc}
}

// Handle godoc
// @Tags journal-config
// @Accept: application/json
// @Produce: application/json
// @Param domain query string false "Domain to check"
// @Param Host header string false "Host header from the request"
// @Success 200 "Domain is allowed"
// @Failure 403 "Domain is not allowed"
// @Router /journal-config/is-allowed [get]
func (this *IsAllowedDomainHandler) Handle(ctx echo.Context) error {

	requestHost := ctx.Request().Header.Get("Host")

	queryDomain := ctx.QueryParam("domain")

	if queryDomain == "" {
		queryDomain = requestHost
	}

	println("[IsAllowedDomainHandler] >>>", queryDomain)
	allowed, err := this.uc.Execute(queryDomain)
	if err != nil {
		return err
	}
	if !allowed {
		return ctx.JSON(403, map[string]string{"error": "Domain is not allowed"})
	}
	return ctx.JSON(200, map[string]string{"message": "Ok"})
}
