package models

import "gorm.io/gorm"

type ArticleAuthorAffiliationModel struct {
	gorm.Model

	ArticleID *uint
	Article   *ArticleModel `gorm:"foreignKey:ArticleID;constraint:OnDelete:SET NULL"`

	AuthorID uint         `gorm:"not null"`
	Author   *AuthorModel `gorm:"foreignKey:AuthorID;constraint:OnDelete:SET NULL"`

	OrganizationID   *uint
	Organization     *OrganizationModel `gorm:"foreignKey:OrganizationID;constraint:OnDelete:SET NULL"`
	OrganizationName string             `gorm:"size:2048"`
	OrganizationTin  string             `gorm:"size:32"`
	PositionName     string             `gorm:"size:255"`
}

func NewArticleAuthorAffiliationModel(articleID *uint, article *ArticleModel, authorID uint, author *AuthorModel, organizationID *uint, organization *OrganizationModel, organizationName string, organizationTin string, positionName string) *ArticleAuthorAffiliationModel {
	return &ArticleAuthorAffiliationModel{ArticleID: articleID, Article: article, AuthorID: authorID, Author: author, OrganizationID: organizationID, Organization: organization, OrganizationName: organizationName, OrganizationTin: organizationTin, PositionName: positionName}
}

func (ArticleAuthorAffiliationModel) TableName() string { return "article_author_affiliations" }
