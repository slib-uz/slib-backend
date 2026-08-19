package gateway

type ORCIDGateway interface {
	AuthorizeUri(redirectUri string) string
	GetByCode(code string, redirectUri string) (string, error)
}
