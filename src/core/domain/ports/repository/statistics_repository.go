package repository

import (
	entity2 "slib.uz/src/core/domain/entity"
)

type StatisticsRepository interface {
	GetSiteStatistics(journalID *uint, publisherID *uint, institutionID *uint) (*entity2.SiteStatisticsEntity, error)
	GetArticleStatisticsByStudyField(page, pageSize int, journalID *uint) (*entity2.PagingEntity[entity2.ArticleStatisticsByStudyFieldEntity], error)
	GetArticleStatisticsByMonth(year, month int, journalID *uint, publisherID *uint) (*entity2.ArticleStatisticsByMonthEntity, error)
	GetArticleStatisticsByLast30Days(journalID *uint, publisherID *uint) (*entity2.ArticleStatisticsByMonthEntity, error)
	GetArticleStatisticsByYear(year int, journalID *uint, publisherID *uint) (*entity2.ArticleStatisticsByYearEntity, error)
}
