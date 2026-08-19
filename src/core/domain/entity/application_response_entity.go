package entity

import (
	"time"
)

type ApplicationResponseEntity struct {
	ID           uint                         `json:"id"`
	Number       string                       `json:"number"`
	ArticleID    uint                         `json:"article_id"`
	Article      *ArticleInputEntity          `json:"article"`
	JournalID    uint                         `json:"journal_id"`
	Journal      *JournalEntity               `json:"journal"`
	UserID       uint                         `json:"user_id"`
	User         *UserEntity                  `json:"user"`
	ReviewStages []*ReviewStageResponseEntity `json:"review_stages"`
	CreatedAt    *time.Time                   `json:"created_at"`
	UpdatedAt    *time.Time                   `json:"updated_at"`
}

func NewApplicationResponseEntity(ID uint, number string, articleID uint, article *ArticleInputEntity, journalID uint, journal *JournalEntity, userID uint, user *UserEntity, reviewStage []*ReviewStageResponseEntity, createdAt *time.Time, updatedAt *time.Time) *ApplicationResponseEntity {
	return &ApplicationResponseEntity{ID: ID, Number: number, ArticleID: articleID, Article: article, JournalID: journalID, Journal: journal, UserID: userID, User: user, ReviewStages: reviewStage, CreatedAt: createdAt, UpdatedAt: updatedAt}
}
