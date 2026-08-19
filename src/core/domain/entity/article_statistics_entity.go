package entity

// ArticleStatisticsByStudyFieldEntity represents statistics grouped by study field
type ArticleStatisticsByStudyFieldEntity struct {
	StudyFieldID   uint              `json:"study_field_id"`
	StudyFieldName map[string]string `json:"study_field_name"`
	Articles       int64             `json:"articles"`
	Journals       int64             `json:"journals"`
}

func NewArticleStatisticsByStudyFieldEntity(studyFieldID uint, studyFieldName map[string]string, articles, journals int64) *ArticleStatisticsByStudyFieldEntity {
	return &ArticleStatisticsByStudyFieldEntity{
		StudyFieldID:   studyFieldID,
		StudyFieldName: studyFieldName,
		Articles:       articles,
		Journals:       journals,
	}
}

// ArticleStatisticsByDayEntity represents statistics for a specific day
type ArticleStatisticsByDayEntity struct {
	Date         string `json:"date"`
	Articles     int64  `json:"articles"`
	Applications int64  `json:"applications"`
}

func NewArticleStatisticsByDayEntity(date string, articles, applications int64) *ArticleStatisticsByDayEntity {
	return &ArticleStatisticsByDayEntity{
		Date:         date,
		Articles:     articles,
		Applications: applications,
	}
}

// ArticleStatisticsByMonthEntity represents statistics for entire month
type ArticleStatisticsByMonthEntity struct {
	Total *ArticleStatisticsTotalEntity            `json:"total"`
	Days  map[string]*ArticleStatisticsByDayEntity `json:"days"`
}

func NewArticleStatisticsByMonthEntity(total *ArticleStatisticsTotalEntity, days map[string]*ArticleStatisticsByDayEntity) *ArticleStatisticsByMonthEntity {
	return &ArticleStatisticsByMonthEntity{
		Total: total,
		Days:  days,
	}
}

// ArticleStatisticsTotalEntity represents total statistics
type ArticleStatisticsTotalEntity struct {
	Articles     int64 `json:"articles"`
	Applications int64 `json:"applications"`
}

func NewArticleStatisticsTotalEntity(articles, applications int64) ArticleStatisticsTotalEntity {
	return ArticleStatisticsTotalEntity{
		Articles:     articles,
		Applications: applications,
	}
}

// ArticleStatisticsByMonthItemEntity represents statistics for a specific month within a year
type ArticleStatisticsByMonthItemEntity struct {
	Articles     int64 `json:"articles"`
	Applications int64 `json:"applications"`
}

func NewArticleStatisticsByMonthItemEntity(articles, applications int64) *ArticleStatisticsByMonthItemEntity {
	return &ArticleStatisticsByMonthItemEntity{
		Articles:     articles,
		Applications: applications,
	}
}

// ArticleStatisticsByYearEntity represents statistics for entire year grouped by month
type ArticleStatisticsByYearEntity struct {
	Total  *ArticleStatisticsTotalEntity                  `json:"total"`
	Months map[string]*ArticleStatisticsByMonthItemEntity `json:"months"`
}

func NewArticleStatisticsByYearEntity(total *ArticleStatisticsTotalEntity, months map[string]*ArticleStatisticsByMonthItemEntity) *ArticleStatisticsByYearEntity {
	return &ArticleStatisticsByYearEntity{
		Total:  total,
		Months: months,
	}
}
