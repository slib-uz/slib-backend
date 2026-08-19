package entity

type UserMeEntity struct {
	ID           uint              `json:"id"`
	FullName     string            `json:"full_name"`
	ScienceID    string            `json:"science_id"`
	Photo        *string           `json:"photo"`
	UserRoles    []*UserRoleEntity `json:"roles"`
	ArticleCount int64             `json:"article_count"`
}

func NewUserMeEntity(id uint, fullName string, scienceID string, photo *string, userRoles []*UserRoleEntity, articleCount int64) *UserMeEntity {
	return &UserMeEntity{ID: id, FullName: fullName, ScienceID: scienceID, Photo: photo, UserRoles: userRoles, ArticleCount: articleCount}
}
