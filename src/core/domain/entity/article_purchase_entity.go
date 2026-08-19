package entity

type ArticlePurchaseEntity struct {
	ID        uint
	JournalID uint
	Journal   *JournalBasicEntity
	ArticleID uint
	Article   *ArticleBasicEntity
	UserID    uint
	User      *UserSharedEntity
	Amount    int64
}

func NewArticlePurchaseEntity(ID uint, journalID uint, journal *JournalBasicEntity, articleID uint, article *ArticleBasicEntity, userID uint, user *UserSharedEntity, amount int64) *ArticlePurchaseEntity {
	return &ArticlePurchaseEntity{ID: ID, JournalID: journalID, Journal: journal, ArticleID: articleID, Article: article, UserID: userID, User: user, Amount: amount}
}

func NewArticlePurchaseCreateEntity(journalID, articleID uint, userID uint, amount int64) *ArticlePurchaseEntity {
	return &ArticlePurchaseEntity{JournalID: journalID, ArticleID: articleID, UserID: userID, Amount: amount}
}
