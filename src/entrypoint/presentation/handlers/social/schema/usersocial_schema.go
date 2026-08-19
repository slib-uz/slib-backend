package schema

type UserSocialResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Link string `json:"link"`
}

type UserSocialRequest struct {
	SocialID uint   `json:"social_id"`
	Link     string `json:"link"`
}
