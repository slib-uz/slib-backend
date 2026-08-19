package entity

type AccessTokenEntity struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"` // seconds
}

func NewAccessTokenEntity(accessToken string, expiresIn int64) *AccessTokenEntity {
	return &AccessTokenEntity{AccessToken: accessToken, ExpiresIn: expiresIn}
}
