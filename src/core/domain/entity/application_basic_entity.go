package entity

type ApplicationBasicEntity struct {
	ID           uint                       `json:"id"`
	Number       string                     `json:"number"`
	ArticleID    uint                       `json:"article_id"`
	Article      *ArticleEntity             `json:"article"`
	JournalID    uint                       `json:"journal_id"`
	Journal      *JournalEntity             `json:"journal"`
	UserID       uint                       `json:"user_id"`
	User         *UserEntity                `json:"user"`
	CurrentStage *ReviewStageResponseEntity `json:"current_stage"`
}

func NewApplicationBasicEntity(ID uint, number string, articleID uint, article *ArticleEntity, journalID uint, journal *JournalEntity, userID uint, user *UserEntity, currentStage *ReviewStageResponseEntity) *ApplicationBasicEntity {
	return &ApplicationBasicEntity{ID: ID, Number: number, ArticleID: articleID, Article: article, JournalID: journalID, Journal: journal, UserID: userID, User: user, CurrentStage: currentStage}
}
