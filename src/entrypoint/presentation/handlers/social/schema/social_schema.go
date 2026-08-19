package schema

type SocialRequest struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type SocialResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}
