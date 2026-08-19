package entity

type ArticleAuthorAffiliationEntity struct {
	ID               uint          `json:"id"`
	ArticleID        *uint         `json:"article_id"`
	AuthorID         uint          `json:"author_id"`
	ScienceID        string        `json:"science_id,omitempty"`
	Author           *AuthorEntity `json:"author"`
	OrganizationID   *uint         `json:"organization_id"`
	OrganizationName string        `json:"organization_name"`
	OrganizationTin  string        `json:"organization_tin"`
	PositionName     string        `json:"position_name"`
}

func NewArticleAuthorAffiliationEntity(id uint, articleID *uint, authorID uint, organizationID *uint, organizationName string, organizationTin string, positionName string) *ArticleAuthorAffiliationEntity {
	return &ArticleAuthorAffiliationEntity{ID: id, ArticleID: articleID, AuthorID: authorID, OrganizationID: organizationID, OrganizationName: organizationName, OrganizationTin: organizationTin, PositionName: positionName}
}
