package schema

type OrcidConnectRequest struct {
	Code        string `json:"code" validate:"required"`
	RedirectUri string `json:"redirect_uri" validate:"required"`
}
