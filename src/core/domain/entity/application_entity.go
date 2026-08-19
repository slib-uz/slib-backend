package entity

import "time"

type ApplicationEntity struct {
	ID           uint                 `json:"id"`
	Number       string               `json:"number"`
	ArticleID    uint                 `json:"article_id"`
	Article      *ArticleEntity       `json:"article,omitempty"`
	JournalID    uint                 `json:"journal_id"`
	Journal      *JournalEntity       `json:"journal,omitempty"`
	UserID       uint                 `json:"user_id"`
	User         *UserEntity          `json:"user,omitempty"`
	IsPublished  bool                 `json:"is_published"`
	ReviewStages []*ReviewStageEntity `json:"review_stages,omitempty"`
	CreatedAt    *time.Time           `json:"-"`
	UpdatedAt    *time.Time           `json:"-"`
}

func NewApplicationEntity(ID uint, number string, articleID uint, article *ArticleEntity, journalID uint, journal *JournalEntity, userID uint, user *UserEntity, isPublished bool, reviewStages []*ReviewStageEntity, createdAt *time.Time, updatedAt *time.Time) *ApplicationEntity {
	return &ApplicationEntity{ID: ID, Number: number, ArticleID: articleID, Article: article, JournalID: journalID, Journal: journal, UserID: userID, User: user, IsPublished: isPublished, ReviewStages: reviewStages, CreatedAt: createdAt, UpdatedAt: updatedAt}
}
