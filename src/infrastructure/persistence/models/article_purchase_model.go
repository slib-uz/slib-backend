package models

import "gorm.io/gorm"

type ArticlePurchaseModel struct {
	gorm.Model

	JournalID uint          `gorm:"index;"`
	Journal   *JournalModel `gorm:"foreignKey:JournalID;references:ID;constraint:OnDelete:SET NULL;"`

	ArticleID uint          `gorm:"uniqueIndex:idx_article_purchase"`
	Article   *ArticleModel `gorm:"foreignKey:ArticleID;references:ID;constraint:OnDelete:CASCADE"`

	UserID uint       `gorm:"not null;uniqueIndex:idx_article_purchase"`
	User   *UserModel `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`

	//InvoiceID uint          `gorm:"not null;"`
	//Invoice   *InvoiceModel `gorm:"foreignKey:InvoiceID;references:ID"`

	Amount int64 `gorm:"not null;default:0"` // Amount of the purchase, can be used for tracking or reporting purposes
}

func NewArticlePurchaseModel(journalID, articleID uint, userID uint, amount int64) *ArticlePurchaseModel {
	return &ArticlePurchaseModel{JournalID: journalID, ArticleID: articleID, UserID: userID, Amount: amount}
}

func (ArticlePurchaseModel) TableName() string {
	return "article_purchases"
}
