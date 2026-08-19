package entity

type UserSocialEntity struct {
	ID            uint   `json:"id"`
	UserProfileID uint   `json:"user_profile_id"`
	Name          string `json:"name"`
	Link          string `json:"link"`
	Icon          string `json:"icon"`
	SocialID      uint   `json:"social_id"`
}

func NewUserSocialEntity(id uint, userProfileID uint, name string, link string, icon string, socialID uint) *UserSocialEntity {
	return &UserSocialEntity{
		ID:            id,
		UserProfileID: userProfileID,
		Name:          name,
		Link:          link,
		Icon:          icon,
		SocialID:      socialID,
	}
}

type UserSocialInputEntity struct {
	ID            uint   `json:"id"`
	UserProfileID uint   `json:"user_profile_id"`
	SocialID      uint   `json:"social_id"`
	Link          string `json:"link"`
}

func NewUserSocialInputEntity(id, userProfileID, socialID uint, link string) *UserSocialInputEntity {
	return &UserSocialInputEntity{
		ID:            id,
		UserProfileID: userProfileID,
		SocialID:      socialID,
		Link:          link,
	}
}
