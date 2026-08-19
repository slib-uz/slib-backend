package repository

import (
	"encoding/json"
	"fmt"

	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
)

type StatisticsRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewStatisticsRepository(baseRepository *BaseRepository) repository.StatisticsRepository {
	return &StatisticsRepositoryImpl{BaseRepository: baseRepository}
}

func (this *StatisticsRepositoryImpl) GetSiteStatistics(journalID *uint, publisherID *uint, institutionID *uint) (*entity2.SiteStatisticsEntity, error) {
	type siteStats struct {
		Users           int64 `json:"users"`
		Articles        int64 `json:"articles"`
		Publishers      int64 `json:"publishers"`
		Journals        int64 `json:"journals"`
		Applications    int64 `json:"applications"`
		Ratings         int64 `json:"ratings"`
		CoAuthors       int64 `json:"co_authors"`
		SupportMessages int64 `json:"support_messages"`
		Comments        int64 `json:"comments"`
	}

	var stats siteStats

	query := `
    SELECT
        (SELECT CASE WHEN ?::bigint IS NOT NULL THEN 0::bigint
            ELSE (SELECT COUNT(*) FROM users WHERE deleted_at IS NULL)::bigint
         END) as users,

        (SELECT COUNT(*) FROM articles a
         LEFT JOIN journals j ON a.journal_id = j.id
         WHERE a.deleted_at IS NULL AND a.is_published = true
         AND (?::bigint IS NULL OR a.journal_id = ?)
         AND (?::bigint IS NULL OR j.publisher_id = ?)
         AND (?::bigint IS NULL OR j.publisher_id IN (
            SELECT p.id FROM publishers p WHERE p.deleted_at IS NULL AND p.institution_id = ?
         ))
        ) as articles,

        (SELECT COUNT(*) FROM publishers p
         WHERE p.deleted_at IS NULL AND p.is_active = true
         AND (?::bigint IS NULL OR p.institution_id = ?)
        ) as publishers,

        (SELECT COUNT(*) FROM journals j
         WHERE j.deleted_at IS NULL AND j.is_active = true AND j.is_accepted = true
         AND (?::bigint IS NULL OR j.publisher_id = ?)
         AND (?::bigint IS NULL OR j.publisher_id IN (
            SELECT p.id FROM publishers p WHERE p.deleted_at IS NULL AND p.institution_id = ?
         ))
        ) as journals,

        (SELECT COUNT(*) FROM article_applications ap
         LEFT JOIN journals j ON ap.journal_id = j.id
         WHERE ap.deleted_at IS NULL
         AND (?::bigint IS NULL OR ap.journal_id = ?)
         AND (?::bigint IS NULL OR j.publisher_id = ?)
         AND (?::bigint IS NULL OR j.publisher_id IN (
            SELECT p.id FROM publishers p WHERE p.deleted_at IS NULL AND p.institution_id = ?
         ))
        ) as applications,

        (SELECT COUNT(*) FROM journal_ratings jr
         INNER JOIN journals j ON jr.journal_id = j.id
         WHERE jr.deleted_at IS NULL
         AND (?::bigint IS NULL OR jr.journal_id = ?)
         AND (?::bigint IS NULL OR j.publisher_id = ?)
         AND (?::bigint IS NULL OR j.publisher_id IN (
            SELECT p.id FROM publishers p WHERE p.deleted_at IS NULL AND p.institution_id = ?
         ))
        ) as ratings,

        (SELECT COUNT(*) FROM comment_models cm
         INNER JOIN articles a ON cm.article_id = a.id
         LEFT JOIN journals j ON a.journal_id = j.id
         WHERE cm.deleted_at IS NULL
         AND (?::bigint IS NULL OR a.journal_id = ?)
         AND (?::bigint IS NULL OR j.publisher_id = ?)
         AND (?::bigint IS NULL OR j.publisher_id IN (
            SELECT p.id FROM publishers p WHERE p.deleted_at IS NULL AND p.institution_id = ?
         ))
        ) as comments,

        (SELECT COUNT(DISTINCT ac.author_model_id)
         FROM article_co_authors ac
         INNER JOIN articles a ON ac.article_model_id = a.id
         LEFT JOIN journals j ON a.journal_id = j.id
         WHERE a.deleted_at IS NULL AND a.is_published = true
         AND (?::bigint IS NULL OR a.journal_id = ?)
         AND (?::bigint IS NULL OR j.publisher_id = ?)
         AND (?::bigint IS NULL OR j.publisher_id IN (
            SELECT p.id FROM publishers p WHERE p.deleted_at IS NULL AND p.institution_id = ?
         ))
        ) as co_authors,

        (SELECT CASE WHEN ?::bigint IS NOT NULL THEN 0::bigint
            ELSE (SELECT COUNT(*) FROM support_dialogs WHERE deleted_at IS NULL AND message_type = 1)::bigint
         END) as support_messages
`

	err := this.db().Raw(query,
		institutionID, // users
		journalID, journalID, publisherID, publisherID, institutionID, institutionID, // articles
		institutionID, institutionID, // publishers
		publisherID, publisherID, institutionID, institutionID, // journals
		journalID, journalID, publisherID, publisherID, institutionID, institutionID, // applications
		journalID, journalID, publisherID, publisherID, institutionID, institutionID, // ratings
		journalID, journalID, publisherID, publisherID, institutionID, institutionID, // comments
		journalID, journalID, publisherID, publisherID, institutionID, institutionID, // co_authors
		institutionID, // support_messages
	).Scan(&stats).Error

	if err != nil {
		return nil, infraError.Wrap(err)
	}

	entity := entity2.NewSiteStatisticsEntity(
		stats.Users,
		stats.Articles,
		stats.Publishers,
		stats.Journals,
		stats.Applications,
		stats.Ratings,
		stats.CoAuthors,
		stats.SupportMessages,
	)
	entity.Ratings = stats.Comments + stats.Ratings
	return entity, nil
}

func (this *StatisticsRepositoryImpl) GetArticleStatisticsByStudyField(page, pageSize int, journalID *uint) (*entity2.PagingEntity[entity2.ArticleStatisticsByStudyFieldEntity], error) {
	type resultRow struct {
		StudyFieldID   uint   `json:"study_field_id"`
		StudyFieldName []byte `json:"study_field_name"` // JSONB as bytes
		Articles       int64  `json:"articles"`
		Journals       int64  `json:"journals"`
	}

	var total int64
	if journalID != nil {
		countQuery := `
            SELECT COUNT(DISTINCT COALESCE(sf.parent_id, sf.id))
            FROM journal_many2many_study_fields jmsf
            JOIN study_fields sf ON jmsf.study_field_model_id = sf.id
            WHERE jmsf.journal_model_id = ?
            AND sf.deleted_at IS NULL
        `
		if err := this.db().Raw(countQuery, *journalID).Scan(&total).Error; err != nil {
			return nil, infraError.Wrap(err)
		}
	} else {
		if err := this.db().Table("study_fields").
			Where("deleted_at IS NULL AND parent_id IS NULL").
			Count(&total).Error; err != nil {
			return nil, infraError.Wrap(err)
		}
	}

	query := `
        WITH parent_study_fields AS (
            SELECT id, name
            FROM study_fields
            WHERE deleted_at IS NULL AND parent_id IS NULL
        ),
        article_stats AS (
            SELECT
                COALESCE(sf.parent_id, sf.id) AS parent_field_id,
                COUNT(DISTINCT a.id) AS article_count
            FROM articles a
            JOIN article_study_fields asf ON a.id = asf.article_model_id
            JOIN study_fields sf ON asf.study_field_model_id = sf.id
            WHERE a.deleted_at IS NULL
            AND a.is_published = TRUE
            AND sf.deleted_at IS NULL
            AND (?::bigint IS NULL OR a.journal_id = ?)
            GROUP BY COALESCE(sf.parent_id, sf.id)
        ),
        journal_stats AS (
            SELECT
                COALESCE(sf.parent_id, sf.id) AS parent_field_id,
                COUNT(DISTINCT j.id) AS journal_count
            FROM journals j
            JOIN journal_many2many_study_fields jmsf ON j.id = jmsf.journal_model_id
            JOIN study_fields sf ON jmsf.study_field_model_id = sf.id
            WHERE j.deleted_at IS NULL
            AND j.is_active = TRUE
            AND j.is_accepted = TRUE
            AND sf.deleted_at IS NULL
            AND (?::bigint IS NULL OR j.id = ?)
            GROUP BY COALESCE(sf.parent_id, sf.id)
        )
        SELECT
            psf.id AS study_field_id,
            psf.name::text AS study_field_name,
            COALESCE(ast.article_count, 0) AS articles,
            COALESCE(jst.journal_count, 0) AS journals
        FROM parent_study_fields psf
        LEFT JOIN article_stats ast ON psf.id = ast.parent_field_id
        LEFT JOIN journal_stats jst ON psf.id = jst.parent_field_id
        WHERE (?::bigint IS NULL OR COALESCE(jst.journal_count, 0) > 0)
        ORDER BY articles DESC, journals DESC
        LIMIT ? OFFSET ?;
    `

	var rawResults []*resultRow
	offset := (page - 1) * pageSize

	if err := this.db().Raw(query,
		journalID, journalID,
		journalID, journalID,
		journalID,
		pageSize, offset,
	).Scan(&rawResults).Error; err != nil {
		return nil, infraError.Wrap(err)
	}

	results := make([]*entity2.ArticleStatisticsByStudyFieldEntity, len(rawResults))
	for i, row := range rawResults {
		var nameMap map[string]string
		if err := json.Unmarshal(row.StudyFieldName, &nameMap); err != nil {
			nameMap = make(map[string]string)
		}
		results[i] = entity2.NewArticleStatisticsByStudyFieldEntity(
			row.StudyFieldID,
			nameMap,
			row.Articles,
			row.Journals,
		)
	}

	return entity2.NewPagingEntity(page, pageSize, total, results), nil
}

func (this *StatisticsRepositoryImpl) GetArticleStatisticsByMonth(year, month int, journalID *uint, publisherID *uint) (*entity2.ArticleStatisticsByMonthEntity, error) {
	type dayStats struct {
		Date         string `json:"date"`
		Articles     int64  `json:"articles"`
		Applications int64  `json:"applications"`
	}

	startDate := fmt.Sprintf("%d-%02d-01", year, month)

	var dbStats []dayStats
	err := this.db().Raw(`
		SELECT
			date_series::date AS date,
			COALESCE(a.articles_count, 0) AS articles,
			COALESCE(ap.applications_count, 0) AS applications
		FROM generate_series(
			?::date,
			(?::date + INTERVAL '1 month - 1 day')::date,
			'1 day'::interval
		) AS date_series
		LEFT JOIN (
			SELECT
				publication_date AS stat_date,
				COUNT(*) AS articles_count
			FROM articles a
			LEFT JOIN journals j ON a.journal_id = j.id
			WHERE a.deleted_at IS NULL
				AND a.is_published = TRUE
				AND a.publication_date >= ?::date
				AND a.publication_date < (?::date + INTERVAL '1 month')::date
				AND (?::bigint IS NULL OR a.journal_id = ?)
				AND (?::bigint IS NULL OR j.publisher_id = ?)
			GROUP BY stat_date
		) a ON date_series::date = a.stat_date
		LEFT JOIN (
			SELECT
				DATE(ap.created_at AT TIME ZONE 'UTC') AS stat_date,
				COUNT(*) AS applications_count
			FROM article_applications ap
			LEFT JOIN journals j ON ap.journal_id = j.id
			WHERE ap.deleted_at IS NULL
				AND ap.created_at >= ?::timestamp
				AND ap.created_at < (?::date + INTERVAL '1 month')::timestamp
				AND (?::bigint IS NULL OR ap.journal_id = ?)
				AND (?::bigint IS NULL OR j.publisher_id = ?)
			GROUP BY stat_date
		) ap ON date_series::date = ap.stat_date
		ORDER BY date_series
	`,
		startDate, startDate, // 1. generate_series
		startDate, startDate, journalID, journalID, publisherID, publisherID, // 2. articles filter
		startDate, startDate, journalID, journalID, publisherID, publisherID, // 3. applications filter
	).Scan(&dbStats).Error
	if err != nil {
		return nil, infraError.Wrap(err)
	}

	days := make(map[string]*entity2.ArticleStatisticsByDayEntity, len(dbStats))
	var totalArticles, totalApplications int64

	for _, stat := range dbStats {
		days[stat.Date] = entity2.NewArticleStatisticsByDayEntity(stat.Date, stat.Articles, stat.Applications)
		totalArticles += stat.Articles
		totalApplications += stat.Applications
	}

	total := entity2.NewArticleStatisticsTotalEntity(totalArticles, totalApplications)
	return entity2.NewArticleStatisticsByMonthEntity(&total, days), nil
}

func (this *StatisticsRepositoryImpl) GetArticleStatisticsByLast30Days(journalID *uint, publisherID *uint) (*entity2.ArticleStatisticsByMonthEntity, error) {
	type dayStats struct {
		Date         string `json:"date"`
		Articles     int64  `json:"articles"`
		Applications int64  `json:"applications"`
	}

	var dbStats []dayStats
	err := this.db().Raw(`
		SELECT
			date_series::date AS date,
			COALESCE(a.articles_count, 0) AS articles,
			COALESCE(ap.applications_count, 0) AS applications
		FROM generate_series(
			(CURRENT_DATE - INTERVAL '29 days')::date,
			CURRENT_DATE::date,
			'1 day'::interval
		) AS date_series
		LEFT JOIN (
			SELECT
				publication_date AS stat_date,
				COUNT(*) AS articles_count
			FROM articles a
			LEFT JOIN journals j ON a.journal_id = j.id
			WHERE a.deleted_at IS NULL
				AND a.is_published = TRUE
				AND a.publication_date >= (CURRENT_DATE - INTERVAL '29 days')::date
				AND a.publication_date <= CURRENT_DATE::date
				AND (?::bigint IS NULL OR a.journal_id = ?)
				AND (?::bigint IS NULL OR j.publisher_id = ?)
			GROUP BY stat_date
		) a ON date_series::date = a.stat_date
		LEFT JOIN (
			SELECT
				DATE(ap.created_at AT TIME ZONE 'UTC') AS stat_date,
				COUNT(*) AS applications_count
			FROM article_applications ap
			LEFT JOIN journals j ON ap.journal_id = j.id
			WHERE ap.deleted_at IS NULL
				AND ap.created_at >= (CURRENT_DATE - INTERVAL '29 days')::timestamp
				AND ap.created_at <= (CURRENT_DATE + INTERVAL '1 day')::timestamp
				AND (?::bigint IS NULL OR ap.journal_id = ?)
				AND (?::bigint IS NULL OR j.publisher_id = ?)
			GROUP BY stat_date
		) ap ON date_series::date = ap.stat_date
		ORDER BY date_series
	`, journalID, journalID, publisherID, publisherID, journalID, journalID, publisherID, publisherID).Scan(&dbStats).Error

	if err != nil {
		return nil, infraError.Wrap(err)
	}

	days := make(map[string]*entity2.ArticleStatisticsByDayEntity, len(dbStats))
	var totalArticles, totalApplications int64

	for _, stat := range dbStats {
		days[stat.Date] = entity2.NewArticleStatisticsByDayEntity(stat.Date, stat.Articles, stat.Applications)
		totalArticles += stat.Articles
		totalApplications += stat.Applications
	}

	total := entity2.NewArticleStatisticsTotalEntity(totalArticles, totalApplications)
	return entity2.NewArticleStatisticsByMonthEntity(&total, days), nil
}

func (this *StatisticsRepositoryImpl) GetArticleStatisticsByYear(year int, journalID *uint, publisherID *uint) (*entity2.ArticleStatisticsByYearEntity, error) {
	type monthStats struct {
		Month        int   `json:"month"`
		Articles     int64 `json:"articles"`
		Applications int64 `json:"applications"`
	}

	var dbStats []monthStats
	err := this.db().Raw(`
		WITH months AS (
			SELECT generate_series(1, 12) AS month
		),
		articles_stats AS (
			SELECT
				EXTRACT(MONTH FROM a.publication_date)::int AS month,
				COUNT(*) AS articles_count
			FROM articles a
			LEFT JOIN journals j ON a.journal_id = j.id
			WHERE a.deleted_at IS NULL
				AND a.is_published = TRUE
				AND EXTRACT(YEAR FROM a.publication_date) = ?
				AND (?::bigint IS NULL OR a.journal_id = ?)
				AND (?::bigint IS NULL OR j.publisher_id = ?)
			GROUP BY EXTRACT(MONTH FROM a.publication_date)
		),
		applications_stats AS (
			SELECT
				EXTRACT(MONTH FROM ap.created_at)::int AS month,
				COUNT(*) AS applications_count
			FROM article_applications ap
			LEFT JOIN journals j ON ap.journal_id = j.id
			WHERE ap.deleted_at IS NULL
				AND EXTRACT(YEAR FROM ap.created_at) = ?
				AND (?::bigint IS NULL OR ap.journal_id = ?)
				AND (?::bigint IS NULL OR j.publisher_id = ?)
			GROUP BY EXTRACT(MONTH FROM ap.created_at)
		)
		SELECT
			m.month,
			COALESCE(a.articles_count, 0) AS articles,
			COALESCE(ap.applications_count, 0) AS applications
		FROM months m
		LEFT JOIN articles_stats a ON a.month = m.month
		LEFT JOIN applications_stats ap ON ap.month = m.month
		ORDER BY m.month
	`, year, journalID, journalID, publisherID, publisherID, year, journalID, journalID, publisherID, publisherID).Scan(&dbStats).Error
	if err != nil {
		return nil, infraError.Wrap(err)
	}

	months := make(map[string]*entity2.ArticleStatisticsByMonthItemEntity, len(dbStats))
	var totalArticles, totalApplications int64

	for _, stat := range dbStats {
		monthKey := fmt.Sprintf("%d", stat.Month)
		months[monthKey] = entity2.NewArticleStatisticsByMonthItemEntity(stat.Articles, stat.Applications)
		totalArticles += stat.Articles
		totalApplications += stat.Applications
	}

	total := entity2.NewArticleStatisticsTotalEntity(totalArticles, totalApplications)
	return entity2.NewArticleStatisticsByYearEntity(&total, months), nil
}
