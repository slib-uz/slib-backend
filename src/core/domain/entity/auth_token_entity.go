package entity

type AuthTokenEntity struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

func NewAuthTokenEntity(accessToken string, refreshToken string) *AuthTokenEntity {
	return &AuthTokenEntity{AccessToken: accessToken, RefreshToken: refreshToken}
}
