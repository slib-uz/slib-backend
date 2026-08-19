package models

import "gorm.io/gorm"

type ArticleApplicationModel struct {
	gorm.Model

	Number string `gorm:"uniqueIndex;size:255"`

	ArticleID uint          `gorm:"index;"`
	Article   *ArticleModel `const:"foreignKey:ArticleID;references:ID;constraint:OnDelete:SET NULL"`

	JournalID uint          `gorm:"index;"`
	Journal   *JournalModel `gorm:"foreignKey:JournalID;references:ID;constraint:OnDelete:SET NULL"`

	UserID uint       `gorm:"index;"`
	User   *UserModel `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:SET NULL"`

	IsPublished bool `gorm:"default:false;not null"`

	ReviewStages      []*ReviewStageModel      `gorm:"foreignKey:ApplicationID;"`
	AntiPlagResults   []*AntiPlagResultModel   `gorm:"foreignKey:ApplicationID;"`
	SpellCheckResults []*SpellCheckResultModel `gorm:"foreignKey:ApplicationID;"`
	AiDetectResults   []*AiDetectResultModel   `gorm:"foreignKey:ApplicationID;"`
}

func (ArticleApplicationModel) TableName() string {
	return "article_applications"
}

func NewArticleApplicationModel(number string, articleID, journalID, userID uint) *ArticleApplicationModel {
	return &ArticleApplicationModel{Number: number, ArticleID: articleID, JournalID: journalID, UserID: userID}
}
