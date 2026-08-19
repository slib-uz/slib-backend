package schema

type UserSocialResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Link string `json:"link"`
}

func NewSocialResponse(id uint, name string, icon string, link string) *UserSocialResponse {
	return &UserSocialResponse{
		ID:   id,
		Name: name,
		Icon: icon,
		Link: link,
	}
}

type UserProfileResponse struct {
	ID        uint                  `json:"id"`
	UserID    uint                  `json:"user_id"`
	FullName  string                `json:"full_name"`
	ScienceID string                `json:"science_id"`
	Email     string                `json:"email"`
	Bio       *string               `json:"bio"`
	Photo     *string               `json:"photo"`
	Phone     string                `json:"phone"`
	Socials   []*UserSocialResponse `json:"socials"`
}

func NewUserProfileResponse(id, userID uint, fullName string, scienceID string, email string, bio *string, photo *string, phone string, socials []*UserSocialResponse) *UserProfileResponse {
	return &UserProfileResponse{
		ID:        id,
		UserID:    userID,
		FullName:  fullName,
		ScienceID: scienceID,
		Email:     email,
		Bio:       bio,
		Photo:     photo,
		Phone:     phone,
		Socials:   socials,
	}
}
