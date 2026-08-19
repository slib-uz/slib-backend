package entity

type MonthStatisticsEntity struct {
	Articles     int64 `json:"articles"`
	Applications int64 `json:"applications"`
}

func NewMonthStatisticsEntity(articles, applications int64) MonthStatisticsEntity {
	return MonthStatisticsEntity{Articles: articles, Applications: applications}
}

type UserArticleStatisticsByYearEntity struct {
	Total  MonthStatisticsEntity            `json:"total"`
	Months map[string]MonthStatisticsEntity `json:"months"`
}

func NewUserArticleStatisticsByYearEntity(total MonthStatisticsEntity, months map[string]MonthStatisticsEntity) *UserArticleStatisticsByYearEntity {
	return &UserArticleStatisticsByYearEntity{Total: total, Months: months}
}

type YearStatisticsEntity struct {
	Articles     int64 `json:"articles"`
	Applications int64 `json:"applications"`
}

func NewYearStatisticsEntity(articles, applications int64) YearStatisticsEntity {
	return YearStatisticsEntity{Articles: articles, Applications: applications}
}

type YearStatisticsItemEntity struct {
	Year         string `json:"year"`
	Articles     int64  `json:"articles"`
	Applications int64  `json:"applications"`
}

func NewYearStatisticsItemEntity(year string, articles, applications int64) YearStatisticsItemEntity {
	return YearStatisticsItemEntity{Year: year, Articles: articles, Applications: applications}
}

type UserArticleStatisticsByYearRangeEntity struct {
	Total MonthStatisticsEntity      `json:"total"`
	Years []YearStatisticsItemEntity `json:"years"`
}

func NewUserArticleStatisticsByYearRangeEntity(total MonthStatisticsEntity, years []YearStatisticsItemEntity) *UserArticleStatisticsByYearRangeEntity {
	return &UserArticleStatisticsByYearRangeEntity{Total: total, Years: years}
}
