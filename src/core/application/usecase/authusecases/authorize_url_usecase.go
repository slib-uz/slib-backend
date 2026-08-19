package authusecases

import (
	"slib.uz/src/core/domain/ports/gateway"
)

type AuthorizeURLUseCase struct {
	gateway gateway.ScienceIdOAuthGateway
}

// NewAuthorizeURLUseCase is a `CONSTRUCTOR` for AuthorizeURLUseCase
// @inject
func NewAuthorizeURLUseCase(gateway gateway.ScienceIdOAuthGateway) *AuthorizeURLUseCase {
	return &AuthorizeURLUseCase{
		gateway: gateway,
	}
}

func (this *AuthorizeURLUseCase) Execute(redirectURL string) string {
	return this.gateway.AuthorizeURL(redirectURL)
}
