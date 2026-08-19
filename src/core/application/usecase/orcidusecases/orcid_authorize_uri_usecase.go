package orcidusecases

import "slib.uz/src/core/domain/ports/gateway"

type OrcidAuthorizeUriUseCase struct {
	gateway gateway.ORCIDGateway
}

// @inject
func NewOrcidAuthorizeUriUseCase(gateway gateway.ORCIDGateway) *OrcidAuthorizeUriUseCase {
	return &OrcidAuthorizeUriUseCase{gateway: gateway}
}

func (this *OrcidAuthorizeUriUseCase) Execute(redirectUrl string) string {
	return this.gateway.AuthorizeUri(redirectUrl)
}
