package entity

type CountStatsEntity struct {
	JournalCount int64 `json:"journal_count"`
	ArticleCount int64 `json:"article_count"`
	AuthorCount  int64 `json:"author_count"`
}

func NewCountStatsEntity(journalCount, articleCount, authorCount int64) *CountStatsEntity {
	return &CountStatsEntity{JournalCount: journalCount, ArticleCount: articleCount, AuthorCount: authorCount}
}
