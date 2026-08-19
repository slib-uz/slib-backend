package schema

type LoginURLRequest struct {
	RedirectURL string `json:"redirect_url" validate:"required,url"`
}

type LoginURLResponse struct {
	RedirectURL string `json:"redirect_url"`
	LoginURL    string `json:"login_url"`
}

func NewLoginURLResponse(redirectURL string, loginURL string) LoginURLResponse {
	return LoginURLResponse{RedirectURL: redirectURL, LoginURL: loginURL}
}
