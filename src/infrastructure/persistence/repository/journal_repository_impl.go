package repository

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/errors"
	_errors "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
	"slib.uz/src/infrastructure/persistence/sorting"
)

type JournalRepositoryImpl struct {
	*BaseRepository
}

func (this *JournalRepositoryImpl) model() *gorm.DB {
	return this.database.GormDB.Model(&models.JournalModel{}).Where("is_accepted", true)
}

// @inject
func NewJournalRepository(repository *BaseRepository) repository.JournalRepository {
	return &JournalRepositoryImpl{BaseRepository: repository}
}

func (this *JournalRepositoryImpl) GetJournalCount() (int64, error) {
	var count int64
	if err := this.model().Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// JournalSortFields — jurnallar ro'yxati uchun ruxsat etilgan tartiblash
// maydonlari. Nomlar journal_list_handler.go dagi Enums(...) ro'yxatiga mos.
var JournalSortFields = sorting.New("journals.established_date DESC", map[string]string{
	"views_count":      "journals.views_count",
	"rating_sum":       "journals.rating_sum",
	"established_date": "journals.established_date",
})

func (this *JournalRepositoryImpl) GetListByPage(page, size int, submissionAccess enum.AccessType, oakRegistered *bool, publisherId uint, name, description, issn, publisherName *string, languageIds, studyFieldIds []uint, fromYear, toYear *int, indexingTypes []enum.IndexingType, sortBy, order string) (*entity.PagingEntity[entity.JournalBasicEntity], error) {
	var journals []*models.JournalModel

	orderExpr, err := JournalSortFields.ResolvePair(sortBy, order)
	if err != nil {
		return nil, err
	}

	query := this.model().
		Preload("Publisher").
		Preload("StudyFields").
		Preload("Indexes").
		Preload("Region").
		Preload("District")

	if submissionAccess != 0 {
		query = query.Where("submission_access = ?", submissionAccess)
	}

	if oakRegistered != nil {
		if *oakRegistered {
			query = query.Where("oak_certificate_file IS NOT NULL")
		} else {
			query = query.Where("oak_certificate_file IS NULL")
		}
	}

	if publisherId > 0 {
		query = query.Where("publisher_id = ?", publisherId)
	}

	if name != nil {
		//query = query.Where("data->>'name' ILIKE ?", "%"+*name+"%")
		query = query.Where(
			"(name->>'uz' ILIKE ? OR name->>'ru' ILIKE ? OR name->>'en' ILIKE ?)",
			"%"+*name+"%", "%"+*name+"%", "%"+*name+"%",
		)
	}

	if description != nil {
		query = query.Where(
			"(description->>'uz' ILIKE ? OR description->>'ru' ILIKE ? OR description->>'en' ILIKE ?)",
			"%"+*description+"%", "%"+*description+"%", "%"+*description+"%",
		)
	}

	if publisherName != nil {
		query = query.
			Joins(`JOIN publishers p ON p.id = journals.publisher_id`).
			Where("(p.name ILIKE ?)", "%"+*publisherName+"%")
	}

	if issn != nil {
		like := "%" + *issn + "%"
		query = query.Where("issn_paper ILIKE ? OR issn_online ILIKE ?", like, like)
	}

	if fromYear != nil {
		query = query.Where("EXTRACT(YEAR FROM established_date) >= ?", *fromYear)
	}
	if toYear != nil {
		query = query.Where("EXTRACT(YEAR FROM established_date) <= ?", *toYear)
	}

	if len(studyFieldIds) > 0 {
		query = query.
			Joins(`JOIN journal_many2many_study_fields jsf ON jsf.journal_model_id = journals.id`).
			Where(`jsf.study_field_model_id IN ?`, studyFieldIds).
			Group("journals.id").
			Having("COUNT(DISTINCT jsf.study_field_model_id) = ?", len(studyFieldIds))

	}

	if len(languageIds) > 0 {
		query = query.Where(`EXISTS (
			SELECT 1 FROM journal_languages jl
			WHERE jl.journal_model_id = journals.id
			AND jl.language_model_id IN ?
		)`, languageIds)
	}

	if len(indexingTypes) > 0 {
		for _, indexingType := range indexingTypes {
			query = query.Where(`EXISTS (
				SELECT 1 FROM journal_indexing ji
				WHERE ji.journal_id = journals.id
				AND ji.indexing_type = ?
			)`, indexingType)
		}
	}

	var count int64
	if result := query.Count(&count); result.Error != nil {
		return nil, errors.Wrap(result.Error)
	}

	query = query.Order(orderExpr)

	if err := query.Limit(size).Offset((page - 1) * size).Find(&journals).Error; err != nil {
		return nil, err
	}

	return mapper.PagingMapper(journals, mapper.JournalModelToBasicEntity, page, size, count), nil
}

func (this *JournalRepositoryImpl) GetByIDWithRelations(id uint) (*entity.JournalEntity, error) {

	var journal models.JournalModel
	if result := this.model().
		Preload("Publisher").
		Preload("StudyFields").
		Preload("Indexes").
		Preload("Languages").
		Preload("Region").
		Preload("District").
		Where("id = ?", id).
		First(&journal); result.Error != nil {
		return nil, errors.Wrap(result.Error)
	}

	var lastID uint
	if result := this.db().Model(&models.JournalApplicationModel{}).
		Where("journal_id = ?", id).
		Order("created_at DESC").
		Limit(1).
		Pluck("id", &lastID); result.Error != nil {
		lastID = 0
	}

	return mapper.JournalModelToEntityWithApplicationID(&journal, lastID), nil
}

func (this *JournalRepositoryImpl) FindByID(id uint) (*entity.JournalEntity, error) {
	var journal models.JournalModel
	if result := this.model().Where("id = ?", id).First(&journal); result.Error != nil {
		return nil, errors.Wrap(result.Error)
	}
	return mapper.JournalModelToEntity(&journal), nil
}

func (this *JournalRepositoryImpl) AddReviewers(journalID uint, reviewerIds []uint) error {

	var _model models.JournalModel

	if err := this.db().First(&_model, journalID).Error; err != nil {
		return errors.Wrap(err)
	}

	reviewers := make([]*models.ReviewerModel, len(reviewerIds))
	for i, reviewerId := range reviewerIds {
		reviewers[i] = &models.ReviewerModel{Model: gorm.Model{ID: reviewerId}}
	}

	err := this.db().Model(&_model).Association("Reviewers").Append(reviewers)

	return errors.Wrap(err)
}

func (repo *JournalRepositoryImpl) UpdateJournal(journalID uint, e *entity.JournalCreateEntity) error {
	return repo.db().Transaction(func(tx *gorm.DB) error {
		var journal models.JournalModel
		if err := tx.First(&journal, journalID).Error; err != nil {
			return err
		}

		updated := mapper.JournalUpdateEntityToModel(&journal, e)

		if err := tx.Model(&journal).
			Select(
				"name", "short_name", "description", "rule", "article_publish_conditions", "address",
				"established_date", "publishing_price", "selling_price", "access_type",
				"issn_paper", "issn_online", "website", "certificate_file", "oak_certificate_file",
				"email", "phone_number", "cover_image_file", "peer_review_method",
				"submission_access", "comment_access", "social_networks", "support_link",
				"region_id", "district_id",
			).
			Updates(&updated).Error; err != nil {
			return err
		}

		if len(e.Indexes) > 0 {
			if err := tx.Where("journal_id = ?", journalID).Delete(&models.JournalIndexingModel{}).Error; err != nil {
				return err
			}
			if err := tx.Create(mapper.JournalIndexingToModel(journalID, e.Indexes)).Error; err != nil {
				return err
			}
		}

		studyFields := make([]*models.StudyFieldModel, len(e.StudyFieldIDs))
		for i, id := range e.StudyFieldIDs {
			studyFields[i] = &models.StudyFieldModel{Model: gorm.Model{ID: id}}
		}
		if err := tx.Model(&journal).Association("StudyFields").Replace(studyFields); err != nil {
			return err
		}

		languages := make([]*models.LanguageModel, len(e.LanguageIDs))
		for i, id := range e.LanguageIDs {
			languages[i] = &models.LanguageModel{Model: gorm.Model{ID: id}}
		}
		return tx.Model(&journal).Association("Languages").Replace(languages)
	})
}

func (this *JournalRepositoryImpl) GetTopJournals(page, pageSize int) (*entity.PagingEntity[entity.JournalEntity], error) {
	var journals []*models.JournalModel
	var total int64

	if err := this.model().
		Count(&total).Error; err != nil {
		return nil, err
	}

	if err := this.model().
		Preload("Publisher").
		Preload("StudyFields").
		Preload("Indexes").
		Order("views_count DESC, rating_sum DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&journals).Error; err != nil {
		return nil, err
	}

	return mapper.PagingMapper(journals, mapper.JournalModelToEntity, page, pageSize, total), nil
}

func (this *JournalRepositoryImpl) UpdateViewsCount(counts map[uint]int64) error {

	if len(counts) == 0 {
		return nil
	}

	ids := make([]uint, 0, len(counts))
	params := make([]any, 0, len(counts)*2)
	caseSQL := "CASE id"

	for id, d := range counts {
		ids = append(ids, id)

		caseSQL += " WHEN ? THEN views_count + ?::bigint"
		params = append(params, id, d)
	}
	caseSQL += " ELSE views_count END"

	return this.db().
		Model(&models.JournalModel{}).
		Where("id IN ?", ids).
		Update("views_count", gorm.Expr(caseSQL, params...)).
		Error
}

func (this *JournalRepositoryImpl) UpdateStatus(journalID uint, isActive bool) error {
	return this.db().Model(&models.JournalModel{}).
		Where("id = ?", journalID).
		Update("is_active", isActive).Error
}

func (this *JournalRepositoryImpl) ExistingIds(ids []uint) ([]uint, error) {
	var existingIds []uint
	if err := this.db().Model(&models.JournalModel{}).Where("is_active = true AND id IN ?", ids).Pluck("id", &existingIds).Error; err != nil {
		return nil, err
	}
	return existingIds, nil
}

func (this *JournalRepositoryImpl) GetPublisherIdByJournalId(journalID uint) (uint, error) {
	var publisherID uint
	if result := this.db().Model(&models.JournalModel{}).Where("id = ?", journalID).Pluck("publisher_id", &publisherID); result.Error != nil {
		return 0, _errors.Wrap(result.Error)
	}
	return publisherID, nil
}

func (this *JournalRepositoryImpl) GetJournalStatistics(page, pageSize int) (*entity.PagingEntity[entity.JournalStatisticEntity], error) {
	var results []struct {
		JournalID              uint
		JournalName            string
		PublisherID            uint
		PublisherName          string
		IssnOnline             string
		IssnPaper              string
		HasStudyFields         bool
		ArticleCount           int64
		JournalPhoneNumber     string
		JournalTelegramContact string
		ChiefEditor            string `gorm:"type:jsonb"`
		Secretary              string `gorm:"type:jsonb"`
	}

	countQuery := `
	SELECT COUNT(*)
	FROM journals j
	WHERE j.is_accepted = true
	`

	var total int64
	if err := this.db().Raw(countQuery).Scan(&total).Error; err != nil {
		return nil, errors.Wrap(err)
	}

	query := `
	SELECT
		j.id as journal_id,
		j.name ->> 'uz' as journal_name,
		p.id as publisher_id,
		p.name as publisher_name,
		coalesce(j.issn_online, '') as issn_online,
		coalesce(j.issn_paper, '') as issn_paper,
		exists (
			select 1
			from journal_many2many_study_fields jsf
			where jsf.journal_model_id = j.id
		) as has_study_fields,
		coalesce(ac.article_count, 0) as article_count,
		coalesce(j.phone_number, '') as journal_phone_number,
		coalesce(j.social_networks ->> 'telegram', '') as journal_telegram_contact,
		rg.chief_editor as chief_editor,
		coalesce(rg.secretary, '[]'::jsonb) as secretary
	FROM journals j
	LEFT JOIN publishers p ON p.id = j.publisher_id
	LEFT JOIN (
		SELECT a.journal_id, count(*) as article_count
		FROM articles a
		WHERE a.deleted_at is null
		GROUP BY a.journal_id
	) ac ON ac.journal_id = j.id
	LEFT JOIN (
	    SELECT
	        r.journal_id,
	        (jsonb_agg(jsonb_build_object(
					'user_id', u.id,
	                'full_name', u.full_name,
	                'science_id', u.science_id,
	                'phone_number', u.phone_number
	        )) FILTER (WHERE r.role = 40))->0 as chief_editor,
	        (jsonb_agg(jsonb_build_object(
					'user_id', u.id,
	                'full_name', u.full_name,
	                'science_id', u.science_id,
	                'phone_number', u.phone_number
	        )) FILTER (WHERE r.role = 50)) as secretary
	    FROM roles r
	             JOIN users u ON u.id = r.user_id AND u.deleted_at is null
	    WHERE r.deleted_at is null AND r.role in (40, 50)
	    GROUP BY r.journal_id
	) rg ON rg.journal_id = j.id
	WHERE j.is_accepted = true
	ORDER BY article_count DESC
	LIMIT ? OFFSET ?
	`

	offset := (page - 1) * pageSize
	if err := this.db().Raw(query, pageSize, offset).Scan(&results).Error; err != nil {
		return nil, errors.Wrap(err)
	}

	statistics := make([]*entity.JournalStatisticEntity, len(results))
	for i, r := range results {
		chiefEditor, _ := entity.UnmarshalJournalStatisticUserEntity(r.ChiefEditor)
		secretaries, _ := entity.UnmarshalJournalStatisticUserEntities(r.Secretary)

		statistics[i] = &entity.JournalStatisticEntity{
			JournalID:              r.JournalID,
			JournalName:            r.JournalName,
			PublisherID:            r.PublisherID,
			PublisherName:          r.PublisherName,
			IssnOnline:             r.IssnOnline,
			IssnPaper:              r.IssnPaper,
			HasStudyFields:         r.HasStudyFields,
			ArticleCount:           r.ArticleCount,
			JournalPhoneNumber:     r.JournalPhoneNumber,
			JournalTelegramContact: r.JournalTelegramContact,
			JournalChiefEditor:     chiefEditor,
			JournalSecretaries:     secretaries,
		}
	}

	return entity.NewPagingEntity(page, pageSize, total, statistics), nil
}

func (this *JournalRepositoryImpl) GetJournalStatisticsV2(page, pageSize int, institutionID, publisherID uint, name, description, issn, publisherName *string) (*entity.PagingEntity[entity.JournalStatisticV2Entity], error) {
	var results []struct {
		JournalID               uint
		JournalName             datatypes.JSON
		PublisherID             uint
		PublisherName           string
		InstitutionID           uint
		InstitutionName         string
		IssnOnline              string
		IssnPaper               string
		JournalRegionID         *uint
		JournalRegionName       *string
		JournalAccessType       enum.AccessType
		JournalSubmissionAccess enum.AccessType
		AdminCount              int64
		EditionCount            int64
		ArticleCount            int64
		ArticleApplicationCount int64
		CoAuthorCount           int64
	}

	filterClause := `
	WHERE j.is_accepted = true AND j.is_active = true
	  AND (? = 0 OR j.publisher_id = ?)
	  AND (? = 0 OR p.institution_id = ?)
	`
	filterArgs := []interface{}{publisherID, publisherID, institutionID, institutionID}

	if name != nil {
		filterClause += ` AND (j.name->>'uz' ILIKE ? OR j.name->>'ru' ILIKE ? OR j.name->>'en' ILIKE ?)`
		like := "%" + *name + "%"
		filterArgs = append(filterArgs, like, like, like)
	}

	if description != nil {
		filterClause += ` AND (j.description->>'uz' ILIKE ? OR j.description->>'ru' ILIKE ? OR j.description->>'en' ILIKE ?)`
		like := "%" + *description + "%"
		filterArgs = append(filterArgs, like, like, like)
	}

	if issn != nil {
		filterClause += ` AND (j.issn_paper ILIKE ? OR j.issn_online ILIKE ?)`
		like := "%" + *issn + "%"
		filterArgs = append(filterArgs, like, like)
	}

	if publisherName != nil {
		filterClause += ` AND (p.name ILIKE ?)`
		filterArgs = append(filterArgs, "%"+*publisherName+"%")
	}

	countQuery := `
	SELECT COUNT(*)
	FROM journals j
	LEFT JOIN publishers p ON p.id = j.publisher_id AND p.deleted_at IS NULL
	` + filterClause

	var total int64
	if err := this.db().Raw(countQuery, filterArgs...).Scan(&total).Error; err != nil {
		return nil, errors.Wrap(err)
	}

	query := `
	SELECT
		j.id as journal_id,
		j.name as journal_name,
		coalesce(p.id, 0) as publisher_id,
		coalesce(p.name, '') as publisher_name,
		coalesce(i.id, 0) as institution_id,
		coalesce(i.name, '') as institution_name,
		coalesce(j.issn_online, '') as issn_online,
		coalesce(j.issn_paper, '') as issn_paper,
		j.region_id as journal_region_id,
		r.name ->> 'uz' as journal_region_name,
		j.access_type as journal_access_type,
		j.submission_access as journal_submission_access,
		coalesce(adm.admin_count, 0) as admin_count,
		coalesce(ec.edition_count, 0) as edition_count,
		coalesce(ac.article_count, 0) as article_count,
		coalesce(aac.article_application_count, 0) as article_application_count,
		coalesce(cac.co_author_count, 0) as co_author_count
	FROM journals j
	LEFT JOIN publishers p ON p.id = j.publisher_id AND p.deleted_at IS NULL
	LEFT JOIN institutions i ON i.id = p.institution_id AND i.deleted_at IS NULL
	LEFT JOIN regions r ON r.id = j.region_id AND r.deleted_at IS NULL
	LEFT JOIN (
		SELECT journal_id, count(*) as admin_count
		FROM roles
		WHERE deleted_at IS NULL AND journal_id IS NOT NULL
		GROUP BY journal_id
	) adm ON adm.journal_id = j.id
	LEFT JOIN (
		SELECT journal_id, count(*) as edition_count
		FROM editions
		WHERE deleted_at IS NULL
		GROUP BY journal_id
	) ec ON ec.journal_id = j.id
	LEFT JOIN (
		SELECT journal_id, count(*) as article_count
		FROM articles
		WHERE deleted_at IS NULL AND is_published = true
		GROUP BY journal_id
	) ac ON ac.journal_id = j.id
	LEFT JOIN (
		SELECT journal_id, count(*) as article_application_count
		FROM article_applications
		WHERE deleted_at IS NULL
		GROUP BY journal_id
	) aac ON aac.journal_id = j.id
	LEFT JOIN (
		SELECT a.journal_id, count(DISTINCT ac.author_model_id) as co_author_count
		FROM article_co_authors ac
		INNER JOIN articles a ON ac.article_model_id = a.id
		WHERE a.deleted_at IS NULL AND a.is_published = true
		GROUP BY a.journal_id
	) cac ON cac.journal_id = j.id
	` + filterClause + `
	ORDER BY article_count DESC
	LIMIT ? OFFSET ?
	`

	offset := (page - 1) * pageSize
	queryArgs := append(append([]interface{}{}, filterArgs...), pageSize, offset)
	if err := this.db().Raw(query, queryArgs...).Scan(&results).Error; err != nil {
		return nil, errors.Wrap(err)
	}

	statistics := make([]*entity.JournalStatisticV2Entity, len(results))
	for i, r := range results {
		statistics[i] = &entity.JournalStatisticV2Entity{
			JournalID:               r.JournalID,
			JournalName:             mapper.FromGormJson[map[string]string](r.JournalName),
			PublisherID:             r.PublisherID,
			PublisherName:           r.PublisherName,
			InstitutionID:           r.InstitutionID,
			InstitutionName:         r.InstitutionName,
			IssnOnline:              r.IssnOnline,
			IssnPaper:               r.IssnPaper,
			JournalRegionID:         r.JournalRegionID,
			JournalRegionName:       r.JournalRegionName,
			JournalAccessType:       r.JournalAccessType,
			JournalSubmissionAccess: r.JournalSubmissionAccess,
			AdminCount:              int(r.AdminCount),
			EditionCount:            int(r.EditionCount),
			ArticleCount:            int(r.ArticleCount),
			ArticleApplicationCount: int(r.ArticleApplicationCount),
			CoAuthorCount:           int(r.CoAuthorCount),
		}
	}

	return entity.NewPagingEntity(page, pageSize, total, statistics), nil
}

func (this *JournalRepositoryImpl) GetJournalStatisticV2ByJournalID(journalID uint) (*entity.JournalStatisticV2Entity, error) {
	var result struct {
		JournalID               uint
		JournalName             datatypes.JSON
		PublisherID             uint
		PublisherName           string
		InstitutionID           uint
		InstitutionName         string
		IssnOnline              string
		IssnPaper               string
		JournalRegionID         *uint
		JournalRegionName       *string
		JournalAccessType       enum.AccessType
		JournalSubmissionAccess enum.AccessType
		AdminCount              int64
		EditionCount            int64
		ArticleCount            int64
		ArticleApplicationCount int64
		CoAuthorCount           int64
	}

	query := `
	SELECT
		j.id as journal_id,
		j.name as journal_name,
		coalesce(p.id, 0) as publisher_id,
		coalesce(p.name, '') as publisher_name,
		coalesce(i.id, 0) as institution_id,
		coalesce(i.name, '') as institution_name,
		coalesce(j.issn_online, '') as issn_online,
		coalesce(j.issn_paper, '') as issn_paper,
		j.region_id as journal_region_id,
		r.name ->> 'uz' as journal_region_name,
		j.access_type as journal_access_type,
		j.submission_access as journal_submission_access,
		coalesce(adm.admin_count, 0) as admin_count,
		coalesce(ec.edition_count, 0) as edition_count,
		coalesce(ac.article_count, 0) as article_count,
		coalesce(aac.article_application_count, 0) as article_application_count,
		coalesce(cac.co_author_count, 0) as co_author_count
	FROM journals j
	LEFT JOIN publishers p ON p.id = j.publisher_id AND p.deleted_at IS NULL
	LEFT JOIN institutions i ON i.id = p.institution_id AND i.deleted_at IS NULL
	LEFT JOIN regions r ON r.id = j.region_id AND r.deleted_at IS NULL
	LEFT JOIN (
		SELECT journal_id, count(*) as admin_count
		FROM roles
		WHERE deleted_at IS NULL AND journal_id IS NOT NULL
		GROUP BY journal_id
	) adm ON adm.journal_id = j.id
	LEFT JOIN (
		SELECT journal_id, count(*) as edition_count
		FROM editions
		WHERE deleted_at IS NULL
		GROUP BY journal_id
	) ec ON ec.journal_id = j.id
	LEFT JOIN (
		SELECT journal_id, count(*) as article_count
		FROM articles
		WHERE deleted_at IS NULL AND is_published = true
		GROUP BY journal_id
	) ac ON ac.journal_id = j.id
	LEFT JOIN (
		SELECT journal_id, count(*) as article_application_count
		FROM article_applications
		WHERE deleted_at IS NULL
		GROUP BY journal_id
	) aac ON aac.journal_id = j.id
	LEFT JOIN (
		SELECT a.journal_id, count(DISTINCT ac.author_model_id) as co_author_count
		FROM article_co_authors ac
		INNER JOIN articles a ON ac.article_model_id = a.id
		WHERE a.deleted_at IS NULL AND a.is_published = true
		GROUP BY a.journal_id
	) cac ON cac.journal_id = j.id
	WHERE j.is_accepted = true AND j.id = ?
	`

	if err := this.db().Raw(query, journalID).Scan(&result).Error; err != nil {
		return nil, errors.Wrap(err)
	}

	if result.JournalID == 0 {
		return nil, nil
	}

	return &entity.JournalStatisticV2Entity{
		JournalID:               result.JournalID,
		JournalName:             mapper.FromGormJson[map[string]string](result.JournalName),
		PublisherID:             result.PublisherID,
		PublisherName:           result.PublisherName,
		InstitutionID:           result.InstitutionID,
		InstitutionName:         result.InstitutionName,
		IssnOnline:              result.IssnOnline,
		IssnPaper:               result.IssnPaper,
		JournalRegionID:         result.JournalRegionID,
		JournalRegionName:       result.JournalRegionName,
		JournalAccessType:       result.JournalAccessType,
		JournalSubmissionAccess: result.JournalSubmissionAccess,
		AdminCount:              int(result.AdminCount),
		EditionCount:            int(result.EditionCount),
		ArticleCount:            int(result.ArticleCount),
		ArticleApplicationCount: int(result.ArticleApplicationCount),
		CoAuthorCount:           int(result.CoAuthorCount),
	}, nil
}

func (this *JournalRepositoryImpl) GetJournalsCompletionStats(journalIds []uint) ([]entity.JournalCompletionStatsEntity, error) {
	if len(journalIds) == 0 {
		return []entity.JournalCompletionStatsEntity{}, nil
	}

	var results []entity.JournalCompletionStatsEntity

	query := `
	SELECT
	    j.id as journal_id,
	    jsonb_exists(j.name, 'uz')        AS has_name,
		jsonb_exists(j.short_name, 'uz') AS has_short_name,
	    (jl.journal_model_id IS NOT NULL)                    AS has_journal_languages,
	    (j.established_date IS NOT NULL)                     AS has_established_date,
	    (j.phone_number IS NOT NULL)                         AS has_phone_number,
	    (j.email IS NOT NULL)                                AS has_email,
	    (j.website IS NOT NULL)                              AS has_website,
	    (j.support_link IS NOT NULL)                         AS has_telegram_support,
	    (j.issn_online IS NOT NULL)                          AS has_issn_online,
	    (j.issn_paper IS NOT NULL)                           AS has_issn_paper,
	    (j.social_networks IS NOT NULL)                      AS has_social_networks,
	    (jsf.journal_model_id IS NOT NULL)                   AS has_study_fields,
	    (COALESCE(j.cover_image_file, '') != '')             AS has_image_file,
	    (j.certificate_file IS NOT NULL)                     AS has_certificate_file,
	    (j.oak_certificate_file IS NOT NULL)                 AS has_oak_certificate_file,
	    (j.address IS NOT NULL)                              AS has_address,
	    (ji.journal_id IS NOT NULL)                          AS has_journal_indexing,
	    COALESCE(LENGTH(j.description ->> 'uz') > 10, FALSE) AS has_description,
	    (j.rule IS NOT NULL)                                 AS has_rule,
	    (j.article_publish_conditions IS NOT NULL)           AS has_article_conditions,
	    (j.region_id IS NOT NULL)                            AS has_region_id

	FROM journals j
	
	LEFT JOIN LATERAL (
	    SELECT journal_model_id
	    FROM journal_languages
	    WHERE journal_model_id = j.id
	    LIMIT 1
	) jl ON TRUE

	LEFT JOIN LATERAL (
	    SELECT journal_model_id
	    FROM journal_many2many_study_fields
	    WHERE journal_model_id = j.id
	    LIMIT 1
	) jsf ON TRUE

	LEFT JOIN LATERAL (
	    SELECT journal_id
	    FROM journal_indexing
	    WHERE journal_id = j.id
	    LIMIT 1
	) ji ON TRUE

	WHERE j.deleted_at IS NULL AND j.id = ANY(?);`

	if err := this.db().Raw(query, pq.Array(journalIds)).Scan(&results).Error; err != nil {
		return nil, errors.Wrap(err)
	}

	return results, nil
}
