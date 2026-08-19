package entity

type SocialEntity struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

func NewSocialEntity(id uint, name string, icon string) *SocialEntity {
	return &SocialEntity{
		ID:   id,
		Name: name,
		Icon: icon,
	}
}
