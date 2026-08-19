package entity

type SiteStatisticsEntity struct {
	Users           int64 `json:"users"`
	Articles        int64 `json:"articles"`
	Publishers      int64 `json:"publishers"`
	Journals        int64 `json:"journals"`
	Applications    int64 `json:"applications"`
	Ratings         int64 `json:"ratings"`
	CoAuthors       int64 `json:"co_authors"`
	SupportMessages int64 `json:"support_messages"`
}

func NewSiteStatisticsEntity(users, articles, publishers, journals, applications, ratings, coAuthors, supportMessages int64) *SiteStatisticsEntity {
	return &SiteStatisticsEntity{
		Users:           users,
		Articles:        articles,
		Publishers:      publishers,
		Journals:        journals,
		Applications:    applications,
		Ratings:         ratings,
		CoAuthors:       coAuthors,
		SupportMessages: supportMessages,
	}
}
