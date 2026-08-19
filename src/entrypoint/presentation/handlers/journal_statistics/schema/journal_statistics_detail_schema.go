package schema

import "slib.uz/src/core/domain/entity"

type JournalStatisticsDetailSchema struct {
	Journal                  entity.JournalEntity `json:"journal"`
	JournalCompletionPercent int                  `json:"journal_completion_percent"`

	ArticleCount                      int `json:"article_count"`
	ArticleApplicationCount           int `json:"article_application_count"`
	CoAuthorCount                     int `json:"coauthor_count"`
	EditionCount                      int `json:"edition_count"`
	Last30DaysArticleApplicationCount int `json:"last_30_days_article_application_count"`
	Last30DaysArticleCount            int `json:"last_30_days_article_count"`

	ArticleStatisticsByYear entity.ArticleStatisticsByYearEntity   `json:"article_statistics_by_year"`
	JournalMembersList      []entity.JournalMemberLastOnlineEntity `json:"journal_members_list"`
}
